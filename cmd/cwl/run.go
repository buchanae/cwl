package main

import (
  "context"
  "os"
  "fmt"
  "encoding/json"
  "path/filepath"
  "github.com/buchanae/cwl"
  "github.com/buchanae/cwl/process"
  localfs "github.com/buchanae/cwl/process/fs/local"
  //gsfs "github.com/buchanae/cwl/process/fs/gs"

  tug "github.com/buchanae/tugboat"
  "github.com/buchanae/tugboat/docker"
  "github.com/buchanae/tugboat/storage/local"
  //gsstore "github.com/buchanae/tugboat/storage/gs"

  "github.com/spf13/cobra"
  "github.com/rs/xid"
)

func init() {
  outdir := "cwl-output"
  debug := false

  cmd := &cobra.Command{
    Use: "run <doc.cwl> <inputs.json>",
    Args: cobra.ExactArgs(2),
    RunE: func(cmd *cobra.Command, args []string) error {
      return run(args[0], args[1], outdir, debug)
    },
  }
  root.AddCommand(cmd)
  f := cmd.Flags()

  f.StringVar(&outdir, "outdir", outdir, "")
  f.BoolVar(&debug, "debug", debug, "")
}

func run(path, inputsPath, outdir string, debug bool) error {
  vals, err := cwl.LoadValuesFile(inputsPath)
  if err != nil {
    return err
  }
  inputsDir := filepath.Dir(inputsPath)

  doc, err := cwl.Load(path)
  if err != nil {
    return err
  }

  r := runner{inputsDir, outdir, debug}

  outvals, err := r.runDoc(doc, vals)
  if err != nil {
    return err
  }

  b, err := json.MarshalIndent(outvals, "", "  ")
  if err != nil {
    return err
  }
  fmt.Println(string(b))

  return err
}

type runner struct {
  inputsDir string
  outdir string
  debug bool
}

func (r *runner) runDoc(doc cwl.Document, vals cwl.Values) (cwl.Values, error) {
  switch z := doc.(type) {
  case *cwl.Workflow:
    return r.runWorkflow(z, vals)
  case *cwl.Tool:
    return r.runTool(z, vals)
  default:
    return nil, fmt.Errorf("running doc: unknown doc type")
  }
}

func (r *runner) runTool(tool *cwl.Tool, vals cwl.Values) (cwl.Values, error) {
  // TODO hack. need to think carefully about how resource requirement and runtime
  //      actually get scheduled.
  var resources *cwl.ResourceRequirement
	reqs := append([]cwl.Requirement{}, tool.Requirements...)
	reqs = append(reqs, tool.Hints...)
  for _, req := range reqs {
    if r, ok := req.(cwl.ResourceRequirement); ok {
      resources = &r
    }
  }

  rt := process.Runtime{}
  // TODO related to the resource requirement search above. basically a hack
  //      for the conformance tests, for now.
  if resources != nil {
    rt.Cores = string(resources.CoresMin)
  }

  fs := localfs.NewLocal(r.inputsDir)
  fs.CalcChecksum = true
  //fs, err := gsfs.NewGS("buchanae-funnel")
  //if err != nil {
    //return nil, err
  //}

  proc, err := process.NewProcess(tool, vals, rt, fs)
  if err != nil {
    return nil, err
  }

  cmd, err := proc.Command()
  if err != nil {
    return nil, err
  }

  //fmt.Fprintln(os.Stderr, cmd)

  workdir := "/cwl"
  // TODO necessary for cwl conformance tests
  image := "python:2"

  if d, ok := tool.RequiresDocker(); ok {
    image = d.Pull
    if d.OutputDirectory != "" {
      workdir = d.OutputDirectory
    }
  }

  task := &tug.Task{
    ID: "cwl-test1-" + xid.New().String(),
    ContainerImage: image,
    Command: cmd,
    Workdir: workdir,
    Volumes: []string{workdir, "/tmp"},
    Env: proc.Env(),

    /* TODO need process.OutputBindings() */
    Outputs: []tug.File{
      {
        URL: r.outdir,
        Path: workdir,
      },
    },
  }
  task.Env["HOME"] = workdir
  task.Env["TMPDIR"] = "/tmp"

  stdout := proc.Stdout()
  stderr := proc.Stderr()
  if stdout != "" {
    task.Stdout = workdir + "/" + stdout
  }
  if stderr != "" {
    task.Stderr = workdir + "/" + stderr
  }

  files := []cwl.File{}
  for _, in := range proc.InputBindings() {
    if f, ok := in.Value.(cwl.File); ok {
      files = append(files, flattenFiles(f)...)
    }
  }
  for _, f := range files {
    task.Inputs = append(task.Inputs, tug.File{
      URL: f.Location,
      // TODO
      Path: f.Path,
    })
  }

  ctx := context.Background()
  store, _ := local.NewLocal()
  //store, _ := gsstore.NewGS("buchanae-funnel")
  var log tug.Logger
  if r.debug {
    log = tug.StderrLogger{}
  } else {
    log = tug.EmptyLogger{}
  }
  exec := &docker.Docker{
    Logger: log,
    NoPull: true,
  }

	stage, err := tug.NewStage("tug-workdir", 0755)
  if err != nil {
    panic(err)
  }
  stage.LeaveDir = true
  defer stage.RemoveAll()

  err = tug.Run(ctx, task, stage, log, store, exec)
  if err != nil {
    if e, ok := err.(*tug.ExecError); ok {
      for _, code := range tool.SuccessCodes {
        if e.ExitCode == code {
          err = nil
        }
      }
    }
  }
  if err != nil {
    return nil, err
  }

  fmt.Fprintln(os.Stderr, "Success")

  //fmt.Println(strings.Join(cmd, " "))

  outfs := localfs.NewLocal(r.outdir)
  outfs.CalcChecksum = true
  //outfs, err := gsfs.NewGS("buchanae-cwl-output")
  if err != nil {
    return nil, err
  }

  return proc.Outputs(outfs)
}

func flattenFiles(file cwl.File) []cwl.File {
  files := []cwl.File{file}
  for _, fd := range file.SecondaryFiles {
    // TODO fix the mismatch between cwl.File and *cwl.File
    if f, ok := fd.(*cwl.File); ok {
      files = append(files, flattenFiles(*f)...)
    }
  }
  return files
}



type Link interface {
  linktype()
  ready() bool
  value() cwl.Value
}

type WorkflowInputLink struct {
  Input cwl.WorkflowInput
  Value cwl.Value
}
func (w *WorkflowInputLink) linktype() {}
func (w *WorkflowInputLink) ready() bool {
  return true
}
func (w *WorkflowInputLink) value() cwl.Value {
  return w.Value
}

type WorkflowStepLink struct {
  Step cwl.Step
  Ready bool
  Value cwl.Value
}
func (w *WorkflowStepLink) linktype() {}
func (w *WorkflowStepLink) ready() bool {
  return w.Ready
}

func (w *WorkflowStepLink) value() cwl.Value {
  return w.Value
}


func (r *runner) runWorkflow(wf *cwl.Workflow, vals cwl.Values) (cwl.Values, error) {

  links := map[string]Link{}

  // TODO input binding
  for _, in := range wf.Inputs {
    val, ok := vals[in.ID]
    if !ok {
      return nil, fmt.Errorf("missing input value for %s", in.ID)
    }
    links[in.ID] = &WorkflowInputLink{in, val}
  }

  for _, step := range wf.Steps {
    for _, out := range step.Out {
      link := &WorkflowStepLink{Step: step}
      links[step.ID +"/"+ out.ID] = link
    }
  }

  // TODO lots of validation, including that all inputs have a valid link.

  remaining := append([]cwl.Step{}, wf.Steps...)

  for len(remaining) > 0 {
    ready, notready := takeReady(remaining, links)
    if len(ready) == 0 {
      break
    }
    remaining = notready

    for _, step := range ready {
      debug("running step", step)
      stepvals := cwl.Values{}

      for _, in := range step.In {
        // TODO handle multiple sources
        if len(in.Source) != 1 {
          panic("multiple sources not implemented")
        }
        src := in.Source[0]
        link := links[src]
        stepvals[in.ID] = link.value()
      }

      // TODO failing because of the cwl.File vs *cwl.File mixup
      debug("STEP VALS", stepvals)
      outvals, err := r.runDoc(step.Run, stepvals)
      if err != nil {
        return nil, err
      }

      for _, out := range step.Out {
        linkID := step.ID +"/"+ out.ID
        link := links[linkID].(*WorkflowStepLink)
        val, ok := outvals[out.ID]
        if !ok {
          return nil, fmt.Errorf("missing output value for %s", linkID)
        }
        link.Value = val
        link.Ready = true
      }
    }
  }

  if len(remaining) > 0 {
    return nil, fmt.Errorf("failed mid-workflow, steps remaining")
  }

  // TODO output values and binding
  return nil, nil
}

func takeReady(steps []cwl.Step, links map[string]Link) (ready, notready []cwl.Step) {
  for _, step := range steps {
    nr := notReady(step.In, links)
    if nr == nil {
      ready = append(ready, step)
    } else {
      notready = append(notready, step)
    }
  }
  return ready, notready
}

// notReady returns step input links which are not ready.
// notReady returns a nil slice if all links are ready.
//
// Note that notReady does not return an error when a link can't be found.
// If a link can't be found, the input will never be ready, so be sure
// to validate links beforehand.
func notReady(inputs []cwl.StepInput, links map[string]Link) []Link {
  var notReady []Link

  for _, in := range inputs {
    for _, src := range in.Source {
      link, ok := links[src]
      if !ok || !link.ready() {
        notReady = append(notReady, link)
      }
    }
  }
  return notReady
}

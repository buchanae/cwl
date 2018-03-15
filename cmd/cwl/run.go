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

  tool, ok := doc.(*cwl.Tool)
  if !ok {
    return fmt.Errorf("can only run command line tools")
  }

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

  fs := localfs.NewLocal(inputsDir)
  fs.CalcChecksum = true
  //fs, err := gsfs.NewGS("buchanae-funnel")
  if err != nil {
    return err
  }

  proc, err := process.NewProcess(tool, vals, rt, fs)
  if err != nil {
    return err
  }

  cmd, err := proc.Command()
  if err != nil {
    return err
  }

  task := &tug.Task{
    ID: "cwl-test1-" + xid.New().String(),
    //ContainerImage: "alpine",
    ContainerImage: "python:2",
    Command: cmd,
    Workdir: "/cwl",
    Volumes: []string{"/cwl", "/tmp"},
    Env: proc.Env(),

    /* TODO need process.OutputBindings() */
    Outputs: []tug.File{
      {
        URL: outdir,
        Path: "/cwl",
      },
    },
  }
  task.Env["HOME"] = "/cwl"
  task.Env["TMPDIR"] = "/tmp"

  stdout, err := proc.Stdout()
  if err != nil {
    return err
  }
  stderr, err := proc.Stderr()
  if err != nil {
    return err
  }
  if stdout != "" {
    task.Stdout = "/cwl/" + stdout
  }
  if stderr != "" {
    task.Stderr = "/cwl/" + stderr
  }

  if d, ok := tool.RequiresDocker(); ok {
    task.ContainerImage = d.Pull
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
  if debug {
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
    return err
  } else {
    fmt.Fprintln(os.Stderr, "Success")
  }

  //fmt.Println(strings.Join(cmd, " "))

  outfs := localfs.NewLocal(outdir)
  outfs.CalcChecksum = true
  //outfs, err := gsfs.NewGS("buchanae-cwl-output")
  if err != nil {
    return err
  }

  outvals, err := proc.Outputs(outfs)
  if err != nil {
    return err
  }

  b, err := json.MarshalIndent(outvals, "", "  ")
  if err != nil {
    return err
  }
  fmt.Println(string(b))

  return nil
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

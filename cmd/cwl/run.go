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
)

func init() {
  outdir := "cwl-output"

  cmd := &cobra.Command{
    Use: "run <doc.cwl> <inputs.json>",
    Args: cobra.ExactArgs(2),
    RunE: func(cmd *cobra.Command, args []string) error {
      return run(args[0], args[1], outdir)
    },
  }
  root.AddCommand(cmd)
  f := cmd.Flags()

  f.StringVar(&outdir, "outdir", outdir, "")
}

func run(path, inputsPath, outdir string) error {
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

  rt := process.Runtime{}
  fs := localfs.NewLocal(inputsDir)
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
    ID: "cwl-test1",
    ContainerImage: "alpine",
    Command: cmd,
    Stdout: "stdout.txt",
    Stderr: "stderr.txt",
    Workdir: "/cwl",
    Volumes: []string{"/cwl"},
    Env: proc.Env(),

    /* TODO need process.OutputBindings() */
    Outputs: []tug.File{
      {
        URL: outdir,
        Path: "/cwl",
      },
      {
        URL: outdir + "/stdout.txt",
        Path: "stdout.txt",
      },
      {
        URL: outdir + "/stderr.txt",
        Path: "stderr.txt",
      },
    },
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
  log := tug.EmptyLogger{}
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
    fmt.Fprintln(os.Stderr, "Error:", err)
  } else {
    fmt.Fprintln(os.Stderr, "Success")
  }

  //fmt.Println(strings.Join(cmd, " "))

  outfs := localfs.NewLocal(outdir)
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

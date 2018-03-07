package main

import (
  "fmt"
  "encoding/json"
  "github.com/buchanae/cwl"
  "github.com/buchanae/cwl/cwllib"
  "github.com/buchanae/cwl/cwllib/env/simple"
  exec "github.com/buchanae/cwl/cwllib/exec/simple"
  "os"
  "github.com/spf13/cobra"
)

var root = cobra.Command{
  Use: "cwl",
	SilenceUsage:  true,
}

func init() {
  cmd := &cobra.Command{
    Use: "dump <doc.cwl>",
    Args: cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
      return dump(args[0])
    },
  }
  root.AddCommand(cmd)
}

func main() {
  if err := root.Execute(); err != nil {
    os.Exit(1)
  }
}

func dump(path string) error {
  doc, err := cwl.Load(path)
  if err != nil {
    return err
  }

  b, err := json.MarshalIndent(doc, "", "  ")
  if err != nil {
    return err
  }

  fmt.Println(string(b))
  return nil
}

func init() {
  cmd := &cobra.Command{
    Use: "run <doc.cwl> <inputs.json>",
    Args: cobra.ExactArgs(2),
    RunE: func(cmd *cobra.Command, args []string) error {
      return run(args[0], args[1])
    },
  }
  root.AddCommand(cmd)
}

func run(path, inputsPath string) error {
  vals, err := cwl.LoadValuesFile(inputsPath)
  if err != nil {
    return err
  }

  doc, err := cwl.Load(path)
  if err != nil {
    return err
  }

  tool, ok := doc.(*cwl.Tool)
  if !ok {
    return fmt.Errorf("can only build command line tools")
  }

  env := simple.NewSimpleEnv()
  job, err := cwllib.NewJob(tool, vals, env)
  if err != nil {
    return err
  }

  cmd, err := job.Command()
  if err != nil {
    return err
  }

  err = exec.Exec(cmd)
  if err != nil {
    return err
  }

  outvals, err := job.Outputs()
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

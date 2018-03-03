package main

import (
  "fmt"
  "encoding/json"
  "github.com/buchanae/cwl"
  "os"
  "github.com/spf13/cobra"
)

var root = cobra.Command{
  Use: "cwl",
	//SilenceErrors: true,
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
    fmt.Println(err)
    os.Exit(1)
  }
}

func dump(path string) error {
  doc, err := cwl.LoadFile(path)
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

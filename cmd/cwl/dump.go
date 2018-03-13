package main

import (
  "fmt"
  "encoding/json"
  "github.com/buchanae/cwl"
  "github.com/spf13/cobra"
)

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

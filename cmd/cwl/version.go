package main

import (
  "fmt"
  "github.com/buchanae/cwl/version"
  "github.com/spf13/cobra"
)

func init() {
  cmd := &cobra.Command{
    Use: "version",
    Args: cobra.NoArgs,
    RunE: func(cmd *cobra.Command, args []string) error {
      fmt.Println(version.String())
      return nil
    },
  }
  root.AddCommand(cmd)
}

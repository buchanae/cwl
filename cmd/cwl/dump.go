package main

import (
  "fmt"
  "encoding/json"
  "github.com/ghodss/yaml"
  "github.com/buchanae/cwl"
  "github.com/spf13/cobra"
)

type dumpOpts struct {
  resolveSchemaDefs bool
  noResolve bool
  json bool
}

func init() {
  opts := dumpOpts{}

  cmd := &cobra.Command{
    Use: "dump <doc.cwl>",
    Args: cobra.ExactArgs(1),
    RunE: func(cmd *cobra.Command, args []string) error {
      return dump(opts, args[0])
    },
  }
  root.AddCommand(cmd)

  f := cmd.Flags()
  f.BoolVar(&opts.resolveSchemaDefs, "resolve-schema-defs", opts.resolveSchemaDefs, "")
  f.BoolVar(&opts.noResolve, "no-resolve", opts.noResolve, "")
  f.BoolVar(&opts.json, "json", opts.json, "")
}

func dump(opts dumpOpts, path string) error {
  var doc cwl.Document
  var err error

  if opts.noResolve {
    doc, err = cwl.LoadWithResolver(path, cwl.NoResolve())
  } else {
    doc, err = cwl.Load(path)
  }
  if err != nil {
    return err
  }

  // TODO also resolve http/file references for schema types?
  if opts.resolveSchemaDefs {
    if tool, ok := doc.(*cwl.Tool); ok {
      err := tool.ResolveSchemaDefs()
      if err != nil {
        return err
      }
    }
  }

  var b []byte
  if opts.json {
    b, err = json.MarshalIndent(doc, "", "  ")
  } else {
    b, err = yaml.Marshal(doc)
  }
  if err != nil {
    return err
  }

  fmt.Println(string(b))
  return nil
}

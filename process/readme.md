This library contains code for processing and executing CWL documents, such as CommandLineTools.

This code is experimental, young, and buggy, however it is able to bind inputs, generate command line arguments, and bind outputs. Currently the focus is on supporting CommandLineTool processing. Last checked (3/14/2018), about 15/68 conformance tests were passing.

## Example

```go
package main

import (
  "fmt"
  "github.com/buchanae/cwl"
  "github.com/buchanae/cwl/process"
  localfs "github.com/buchanae/cwl/process/fs/local"
)

func main() {
  // must be run from /examples dir
  path := "tar-param.cwl"
  inputsPath := "tar-param.inputs.yml"

  vals, err := cwl.LoadValuesFile(inputsPath)
  if err != nil {
    fmt.Println("error loading inputs:", err)
    return
  }

  doc, err := cwl.Load(path)
  if err != nil {
    panic(err)
  }

  tool, ok := doc.(*cwl.Tool)
  if !ok {
    panic("can only run command line tools")
  }

  rt := process.Runtime{}
  fs := localfs.NewLocal(".")

  proc, err := process.NewProcess(tool, vals, rt, fs)
  if err != nil {
    panic(err)
  }

  cmd, err := proc.Command()
  if err != nil {
    panic(err)
  }
  fmt.Println(cmd)
}
```

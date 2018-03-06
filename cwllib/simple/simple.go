package simple

import (
  "github.com/buchanae/cwl/cwllib"
  "os"
  "os/exec"
)

// TODO
// working directory, output directory
// output file reporting
// environment
// async handle
// resource matching
// exit code checking
//
// check for cwl.output.json

type Result struct {
  ExitCode int
}

func Exec(job *cwllib.Job) error {
  cmd := exec.Command(job.Command[0], job.Command[1:]...)
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  return cmd.Run()
}

package simple

import (
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

func Exec(args []string) error {
  cmd := exec.Command(args[0], args[1:]...)
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr
  return cmd.Run()
}

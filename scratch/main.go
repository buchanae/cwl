package main

import (
  "os"
  "github.com/ohsu-comp-bio/cwl"
)

func main() {
  p := os.Args[1]
  cwl.LoadFile(p)
}

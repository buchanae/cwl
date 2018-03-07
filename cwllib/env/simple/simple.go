package simple

import (
  "github.com/buchanae/cwl"
  "github.com/buchanae/cwl/cwllib"
  "github.com/buchanae/cwl/cwllib/fs/local"
)

type SimpleEnv struct {
  fs *local.Local
}

func NewSimpleEnv() *SimpleEnv {
	/*
	  id, err := uuid.NewRandom()
	  if err != nil {
	    return errf("error generating unique file location: %s", err)
	  }
	*/
  return &SimpleEnv{fs: local.NewLocal(".")}
}

func (s *SimpleEnv) Runtime() cwllib.Runtime {
  return cwllib.Runtime{}
}

func (s *SimpleEnv) Filesystem() cwllib.Filesystem {
  return s.fs
}

func (s *SimpleEnv) SupportsDocker() bool {
  return false
}

func (s *SimpleEnv) SupportsShell() bool {
  return false
}

func (s *SimpleEnv) CheckResources(req cwl.ResourceRequirement) error {
  return nil
}

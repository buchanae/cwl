package cwllib

import (
	"github.com/buchanae/cwl"
)

type Env interface {
  Runtime() Runtime
  Filesystem() Filesystem
}

type Mebibyte int

type Runtime struct {
  Outdir string
  Tmpdir string
  Cores int
  RAM Mebibyte
  OutdirSize Mebibyte
  TmpdirSize Mebibyte
}

type Job struct {
  tool *cwl.Tool
  inputs cwl.Values
  env Env
  bindings []*binding
  expressionLibs []string
}

func NewJob(tool *cwl.Tool, inputs cwl.Values, env Env) (*Job, error) {

  err := ValidateTool(tool)
  if err != nil {
    return nil, err
  }

  // TODO expose input bindings as an exported type of data
  //      could be useful to know separately from all the other processing.
  job := &Job{
    tool: tool,
    inputs: inputs,
    env: env,
  }

	// Bind inputs to values.
  //
  // Since every part of a tool depends on "inputs" being available to expressions,
  // nothing can be done on a Job without a valid inputs binding,
  // which is why we bind in the Job constructor.
	for i, in := range tool.Inputs {
		val := inputs[in.ID]
		if val == nil {
			val = in.Default
		}

		k := sortKey{getPos(in.InputBinding), i}
		b, err := job.bindInput(in.Type, in.InputBinding, in.SecondaryFiles, val, k)
		if err != nil {
			return nil, errf("error while binding inputs: %s", err)
		}
		if b == nil {
			return nil, errf("no binding found for input: %s", in.ID)
		}

    job.bindings = append(job.bindings, b...)
	}

  return job, nil
}

func (job *Job) eval(expr cwl.Expression, self interface{}) (interface{}, error) {
  data := ExprData{
    Inputs: job.inputs,
    Self: self,
    Libs: job.expressionLibs,
    Runtime: job.env.Runtime(),
  }
  return Eval(expr, data)
}

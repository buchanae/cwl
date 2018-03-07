package cwllib

import (
	"github.com/buchanae/cwl"
)

type Env interface {
	Runtime() Runtime
	Filesystem() Filesystem
	SupportsDocker() bool
	SupportsShell() bool
	CheckResources(cwl.ResourceRequirement) error
}

type Mebibyte int

// TODO this is provided to expressions early on in job processing,
//      but it won't have real values from a scheduler until much later.
type Runtime struct {
	Outdir     string
	Tmpdir     string
	Cores      int
	RAM        Mebibyte
	OutdirSize Mebibyte
	TmpdirSize Mebibyte
}

type Job struct {
	tool           *cwl.Tool
	inputs         cwl.Values
	env            Env
	bindings       []*binding
	expressionLibs []string
	envvars        map[string]string
	docker         cwl.DockerRequirement
	resources      cwl.ResourceRequirement
	shell          bool
}

func NewJob(tool *cwl.Tool, inputs cwl.Values, env Env) (*Job, error) {

	err := ValidateTool(tool)
	if err != nil {
		return nil, err
	}

	// TODO expose input bindings as an exported type of data
	//      could be useful to know separately from all the other processing.
	job := &Job{
		tool:    tool,
		inputs:  inputs,
		env:     env,
		envvars: map[string]string{},
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

	err = job.loadReqs(tool.Requirements)
	if err != nil {
		return nil, err
	}
	err = job.loadReqs(tool.Hints)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (job *Job) loadReqs(reqs []cwl.Requirement) error {
	for _, req := range reqs {
		switch z := req.(type) {

		case cwl.EnvVarRequirement:
			err := job.evalEnvVars(z.EnvDef)
			if err != nil {
				return errf("failed to evaluate EnvVarRequirement: %s", err)
			}

		case cwl.InlineJavascriptRequirement:
			job.expressionLibs = append(job.expressionLibs, z.ExpressionLib...)

		case cwl.DockerRequirement:
			if !job.env.SupportsDocker() {
				return errf("The selected compute environment does not support Docker")
			}
			job.docker = z

		case cwl.ResourceRequirement:
			err := job.env.CheckResources(z)
			if err != nil {
				return errf("ResourceRequirement failed: %s", err)
			}
			job.resources = z

		case cwl.ShellCommandRequirement:
			if !job.env.SupportsShell() {
				return errf("The selected compute environment does not support shell commands")
			}
			job.shell = true

		case cwl.SchemaDefRequirement:
			return errf("SchemaDefRequirement is not supported (yet)")
		case cwl.InitialWorkDirRequirement:
			return errf("InitialWorkDirRequirement is not supported (yet)")
		default:
			return errf("unknown requirement type")
		}
	}
	return nil
}

func (job *Job) evalEnvVars(def map[string]cwl.Expression) error {
	for k, expr := range def {
		val, err := job.eval(expr, nil)
		if err != nil {
			return errf(`failed to evaluate expression: "%s": %s`, expr, err)
		}
		str, ok := val.(string)
		if !ok {
			return errf(`EnvVar must evaluate to a string, got "%s"`, val)
		}
		job.envvars[k] = str
	}
	return nil
}

func (job *Job) eval(expr cwl.Expression, self interface{}) (interface{}, error) {
	data := ExprData{
		Inputs:  job.inputs,
		Self:    self,
		Libs:    job.expressionLibs,
		Runtime: job.env.Runtime(),
	}
	return Eval(expr, data)
}

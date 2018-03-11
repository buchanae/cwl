package process

import (
	"github.com/buchanae/cwl"
	"github.com/buchanae/cwl/expr"
)

type Env interface {
	Runtime() Runtime
	Filesystem() Filesystem
	CheckResources(cwl.ResourceRequirement) error
}

type Mebibyte int

// TODO this is provided to expressions early on in process processing,
//      but it won't have real values from a scheduler until much later.
type Runtime struct {
	Outdir     string
	Tmpdir     string
	Cores      int
	RAM        Mebibyte
	OutdirSize Mebibyte
	TmpdirSize Mebibyte
}

type Process struct {
	tool           *cwl.Tool
	inputs         cwl.Values
	env            Env
	bindings       []*binding
	expressionLibs []string
	envvars        map[string]string
	resources      cwl.ResourceRequirement
	shell          bool
}

func NewProcess(tool *cwl.Tool, inputs cwl.Values, env Env) (*Process, error) {

	err := cwl.ValidateTool(tool)
	if err != nil {
		return nil, err
	}

	// TODO expose input bindings as an exported type of data
	//      could be useful to know separately from all the other processing.
	process := &Process{
		tool:    tool,
		inputs:  inputs,
		env:     env,
		envvars: map[string]string{},
	}

	// Bind inputs to values.
	//
	// Since every part of a tool depends on "inputs" being available to expressions,
	// nothing can be done on a Process without a valid inputs binding,
	// which is why we bind in the Process constructor.
	for i, in := range tool.Inputs {
		val := inputs[in.ID]
		if val == nil {
			val = in.Default
		}

		k := sortKey{getPos(in.InputBinding), i}
		b, err := process.bindInput(in.Type, in.InputBinding, in.SecondaryFiles, val, k)
		if err != nil {
			return nil, errf("error while binding inputs: %s", err)
		}
		if b == nil {
			return nil, errf("no binding found for input: %s", in.ID)
		}

		process.bindings = append(process.bindings, b...)
	}

	if libs, ok := tool.RequiresInlineJavascript(); ok {
		process.expressionLibs = libs
	}

	process.shell = tool.RequiresShellCommand()

	err = process.loadReqs(tool.Requirements)
	if err != nil {
		return nil, err
	}
	err = process.loadReqs(tool.Hints)
	if err != nil {
		return nil, err
	}

	return process, nil
}

func (process *Process) loadReqs(reqs []cwl.Requirement) error {
	for _, req := range reqs {
		switch z := req.(type) {

		case cwl.EnvVarRequirement:
			err := process.evalEnvVars(z.EnvDef)
			if err != nil {
				return errf("failed to evaluate EnvVarRequirement: %s", err)
			}

		case cwl.ResourceRequirement:
			err := process.env.CheckResources(z)
			if err != nil {
				return errf("ResourceRequirement failed: %s", err)
			}
			process.resources = z

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

func (process *Process) evalEnvVars(def map[string]cwl.Expression) error {
	for k, expr := range def {
		val, err := process.eval(expr, nil)
		if err != nil {
			return errf(`failed to evaluate expression: "%s": %s`, expr, err)
		}
		str, ok := val.(string)
		if !ok {
			return errf(`EnvVar must evaluate to a string, got "%s"`, val)
		}
		process.envvars[k] = str
	}
	return nil
}

func (process *Process) eval(x cwl.Expression, self interface{}) (interface{}, error) {
	r := process.env.Runtime()
	return expr.Eval(x, process.expressionLibs, map[string]interface{}{
		"inputs": process.inputs,
		"self":   self,
		"runtime": map[string]interface{}{
			"outdir":     r.Outdir,
			"tmpdir":     r.Tmpdir,
			"cores":      r.Cores,
			"ram":        r.RAM,
			"outdirSize": r.OutdirSize,
			"tmpdirSize": r.TmpdirSize,
		},
	})
}

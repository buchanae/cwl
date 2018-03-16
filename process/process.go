package process

import (
	"github.com/buchanae/cwl"
	"github.com/buchanae/cwl/expr"
	"github.com/rs/xid"
)

type Mebibyte int

// TODO this is provided to expressions early on in process processing,
//      but it won't have real values from a scheduler until much later.
type Runtime struct {
	Outdir string
	Tmpdir string
	// TODO make these all strings?
	Cores      string
	RAM        Mebibyte
	OutdirSize Mebibyte
	TmpdirSize Mebibyte
}

type Resources struct {
	CoresMin,
	CoresMax int

	RAMMin,
	RAMMax,
	OutdirMin,
	OutdirMax,
	TmpdirMin,
	TmpdirMax Mebibyte
}

type Process struct {
	tool           *cwl.Tool
	inputs         cwl.Values
	runtime        Runtime
	fs             Filesystem
	bindings       []*Binding
	expressionLibs []string
	env            map[string]string
	shell          bool
	resources      Resources
	stdout         string
	stderr         string
}

func NewProcess(tool *cwl.Tool, inputs cwl.Values, rt Runtime, fs Filesystem) (*Process, error) {

	err := cwl.ValidateTool(tool)
	if err != nil {
		return nil, err
	}

	// TODO expose input bindings as an exported type of data
	//      could be useful to know separately from all the other processing.
	process := &Process{
		tool:    tool,
		inputs:  inputs,
		runtime: rt,
		fs:      fs,
		env:     map[string]string{},
	}

	// Bind inputs to values.
	//
	// Since every part of a tool depends on "inputs" being available to expressions,
	// nothing can be done on a Process without a valid inputs binding,
	// which is why we bind in the Process constructor.
	for _, in := range tool.Inputs {
		val := inputs[in.ID]
		if val == nil {
			val = in.Default
		}

		k := sortKey{getPos(in.InputBinding)}
		b, err := process.bindInput(in.ID, in.Type, in.InputBinding, in.SecondaryFiles, val, k)
		if err != nil {
			return nil, errf("binding input %q: %s", in.ID, err)
		}
		if b == nil {
			return nil, errf("no binding found for input: %s", in.ID)
		}

		process.bindings = append(process.bindings, b...)
	}

	err = process.loadReqs()
	if err != nil {
		return nil, err
	}

	stdoutI, err := process.eval(process.tool.Stdout, nil)
	if err != nil {
		return nil, wrap(err, "evaluating stdout expression")
	}

	stderrI, err := process.eval(process.tool.Stderr, nil)
	if err != nil {
		return nil, wrap(err, "evaluating stderr expression")
	}

	var stdoutStr, stderrStr string
	var ok bool

	if stdoutI != nil {
		stdoutStr, ok = stdoutI.(string)
		if !ok {
			return nil, errf("stdout expression returned a non-string value")
		}
	}

	if stderrI != nil {
		stderrStr, ok = stderrI.(string)
		if !ok {
			return nil, errf("stderr expression returned a non-string value")
		}
	}

	for _, out := range process.tool.Outputs {
		if len(out.Type) == 1 {
			if _, ok := out.Type[0].(cwl.Stdout); ok && stdoutStr == "" {
				stdoutStr = "stdout-" + xid.New().String()
			}
			if _, ok := out.Type[0].(cwl.Stderr); ok && stderrStr == "" {
				stderrStr = "stderr-" + xid.New().String()
			}
		}
	}
	process.stdout = stdoutStr
	process.stderr = stderrStr

	return process, nil
}

func (process *Process) Stdout() string {
	return process.stdout
}

func (process *Process) Stderr() string {
	return process.stderr
}

func (process *Process) Tool() *cwl.Tool {
	return process.tool
}

func (process *Process) Resources() Resources {
	return process.resources
}

func (process *Process) InputBindings() []*Binding {
	bindings := make([]*Binding, len(process.bindings))
	copy(bindings, process.bindings)
	return bindings
}

func (process *Process) Env() map[string]string {
	env := map[string]string{}
	for k, v := range process.env {
		env[k] = v
	}
	return env
}

func (process *Process) loadReqs() error {
	reqs := append([]cwl.Requirement{}, process.tool.Requirements...)
	reqs = append(reqs, process.tool.Hints...)

	for _, req := range reqs {
		switch z := req.(type) {

		case cwl.InlineJavascriptRequirement:
			process.expressionLibs = z.ExpressionLib

		case cwl.EnvVarRequirement:
			err := process.evalEnvVars(z.EnvDef)
			if err != nil {
				return errf("failed to evaluate EnvVarRequirement: %s", err)
			}

		case cwl.ResourceRequirement:
			// TODO eval expressions

		case cwl.SchemaDefRequirement:
			return errf("SchemaDefRequirement is not supported (yet)")
		case cwl.InitialWorkDirRequirement:
			return errf("InitialWorkDirRequirement is not supported (yet)")
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
		process.env[k] = str
	}
	return nil
}

func (process *Process) eval(x cwl.Expression, self interface{}) (interface{}, error) {
	r := process.runtime
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

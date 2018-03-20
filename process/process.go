package process

import (
	"encoding/json"
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

func NewProcess(tool *cwl.Tool, values cwl.Values, rt Runtime, fs Filesystem) (*Process, error) {

	err := cwl.ValidateTool(tool)
	if err != nil {
		return nil, err
	}

	// TODO expose input bindings as an exported type of data
	//      could be useful to know separately from all the other processing.
	process := &Process{
		tool:    tool,
		inputs:  values,
		runtime: rt,
		fs:      fs,
		env:     map[string]string{},
	}

	// Set default input values.
	setDefaults(values, tool.Inputs)

	// Bind inputs to values.
	//
	// Since every part of a tool depends on "inputs" being available to expressions,
	// nothing can be done on a Process without a valid inputs binding,
	// which is why we bind in the Process constructor.
	for _, in := range tool.Inputs {
		val := values[in.ID]
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
	// TODO copying slice, but still using pointers. deep copy?
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

	inputsData := map[string]interface{}{}
	for _, b := range process.bindings {
		v, err := toJSONMap(b.Value)
		if err != nil {
			return nil, wrap(err, `mashaling "%s" for JS eval`, b.name)
		}
		if v == nil {
			v = expr.Null
		}
		inputsData[b.name] = v
	}

	selfData, err := toJSONMap(self)
	if err != nil {
		return nil, wrap(err, `marshaling "self" for JS eval`)
	}

	r := process.runtime
	return expr.Eval(x, process.expressionLibs, map[string]interface{}{
		"inputs": inputsData,
		"self":   selfData,
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

func toJSONMap(v interface{}) (interface{}, error) {
	if v == nil {
		return nil, nil
	}

	// Need to convert Go variable naming to JSON. Easiest way to to marshal to JSON,
	// then unmarshal into a map.
	j, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var data interface{}
	err = json.Unmarshal(j, &data)
	if err != nil {
		return nil, wrap(err, `marshaling data for JS evaluation`)
	}
	return data, nil
}

// setDefaults sets the default input values based on the CommandInput.Default.
func setDefaults(values cwl.Values, inputs []cwl.CommandInput) {
	for _, in := range inputs {
		_, ok := values[in.ID]
		if !ok && in.Default != nil {
			values[in.ID] = in.Default
		}
	}
}

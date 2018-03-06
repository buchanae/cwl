package cwllib

import (
	"fmt"
	"github.com/buchanae/cwl"
	"github.com/spf13/cast"
	"sort"
)

/*
TODO
- better exec vs document code organization
- more complete JS expression context (self, inputs, runtime, etc)
- output document binding
- cwl.output.json
- secondary files
- load expression result values into File/Directory types where appropriate
- file staging and working directory
- solid expression parser (regexp misses edge cases and escaping)
- relative path context (current working directory) for filesystems
- absolute paths for files, especially in outputs
- resolve document references
- filesystem multiplexing based on location
- success/failure codes and relationship to CLI cmd
- Any type
- document validation before processing
- better line/col info from document loading errors
- carefully check document json/yaml marshaling
- input/output record type handling
- executor backends
- directory type
- good framework for e2e tests with lots of coverage
- $include and $import
- test unrecognized fields are ignored (possibly with warning)
- optional checksum calculation for filesystems
- resource requests
- environment variables
- initial work dir
- docker
- missing requirement/hint types. see requirements.go

workflow execution:
- basics
- caching

server + API:
*/

type Job struct {
	Command []string
}

type Executor struct {
	FS Filesystem
}

func NewExecutor() *Executor {
	return &Executor{FS: NewLocal()}
}

// evalGlobPatterns evaluates a list of potential expressions as defined by the CWL
// OutputBinding.glob field. It returns a list of strings, which are glob expression,
// or an error.
//
// cwl spec:
// "If an expression is provided, the expression must return a string or an array
//  of strings, which will then be evaluated as one or more glob patterns."
func evalGlobPatterns(patterns []cwl.Expression) ([]string, error) {
	var out []string

	for _, pattern := range patterns {
		val, err := Eval(pattern)
		if err != nil {
			return nil, err
		}

		switch z := val.(type) {
		case string:
			out = append(out, z)
		case []string:
			out = append(out, z...)
		default:
			return nil, fmt.Errorf(
				"glob expression returned an invalid type. Only string or []string "+
					"are allowed. Got: %s", val)
		}
	}
	return out, nil
}

// matchFiles executes the list of glob patterns, returning a list of matched files.
// matchFiles must return a non-nil list on success, even if no files are matched.
func (e *Executor) matchFiles(globs []string, loadContents bool) ([]*cwl.File, error) {
	// it's important this slice isn't nil, because the outputEval field
	// expects it to be non-null during expression evaluation.
	files := []*cwl.File{}

	// resolve all the globs into file objects.
	for _, pattern := range globs {
		matches, err := e.FS.Glob(pattern)
		if err != nil {
			return nil, fmt.Errorf("failed to execute glob: %s", err)
		}

		for _, m := range matches {
			// TODO clean this up. the split between this and the "fs" package is weird.
			v := cwl.File{
				Location: m.Location,
				Path:     m.Path,
				Checksum: m.Checksum,
				Size:     m.Size,
			}

			f, err := ResolveFile(v, e.FS, loadContents)
			if err != nil {
				return nil, err
			}
			files = append(files, f)
		}
	}
	return files, nil
}

// TODO the expressions here need access to "inputs"
// CollectOutputs collects outputs from the given CommandLineTool.
func (e *Executor) CollectOutputs(clt *cwl.CommandLineTool) (cwl.Values, error) {
	values := cwl.Values{}
	for _, out := range clt.Outputs {
		v, err := e.CollectOutput(out)
		if err != nil {
			return nil, err
		}
		values[out.ID] = v
	}
	return values, nil
}

// CollectOutput collects the output value for a single CommandOutput.
func (e *Executor) CollectOutput(out cwl.CommandOutput) (val interface{}, err error) {
	// glob patterns may be expressions. evaluate them.
	globs, err := evalGlobPatterns(out.OutputBinding.Glob)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate glob expressions for %s: %s", out.ID, err)
	}

	files, err := e.matchFiles(globs, out.OutputBinding.LoadContents)
	if err != nil {
		return nil, fmt.Errorf("failed to match files for %s: %s", out.ID, err)
	}

	if out.OutputBinding.OutputEval != "" {
		// TODO set value of "self"
		val, err := Eval(out.OutputBinding.OutputEval)
		if err != nil {
			return nil, fmt.Errorf("failed to evaluate outputEval for %s: %s", out.ID, err)
		}
		return val, nil
	}
	return files, nil
}

// BuildJob builds command line arguments for an invocation a tool
// given a set of input values.
func (e *Executor) BuildJob(clt *cwl.CommandLineTool, vals cwl.InputValues) (*Job, error) {
	args := bindings{}

	// Add "arguments"
	for i, arg := range clt.Arguments {
		// TODO validate that valueFrom is set
		val, err := Eval(arg.ValueFrom)
		if err != nil {
			return nil, fmt.Errorf("failed to eval argument value: %s", err)
		}
		b := &binding{arg, argType{}, val, sortKey{arg.Position, i}, nil}
		args = append(args, b)
	}

	// Bind inputs to values and add args.
	for i, in := range clt.Inputs {
		k := sortKey{getPos(in.InputBinding), i}
		val := vals[in.ID]
		if val == nil {
			val = in.Default
		}

		b, err := e.walk(in, val, k)
		if err != nil {
			return nil, fmt.Errorf("error while binding inputs: %s", err)
		}
		if b == nil {
			return nil, fmt.Errorf("no binding found for input: %s", in.ID)
		}

		args = append(args, b...)
	}

	sort.Stable(args)
	//debug(args)

	job := &Job{
		Command: append([]string{}, clt.BaseCommand...),
	}

	// Now collect the input bindings into command line arguments
	for _, b := range args {
		job.Command = append(job.Command, b.args()...)
	}

	return job, nil
}

// walk walks the tree of input descriptors and values,
// binding values to descriptors.
//
// walk is called recursively for types which have subtypes,
// such as array, record, etc.
func (e *Executor) walk(b interface{}, val interface{}, key sortKey) (bindings, error) {
	var types []cwl.InputType
	var clb *cwl.CommandLineBinding

	switch z := b.(type) {
	case cwl.CommandInput:
		types, clb = z.Type, z.InputBinding
	case cwl.InputArray:
		types, clb = z.Items, z.InputBinding
	case cwl.InputField:
		types, clb = z.Type, z.InputBinding
	default:
		panic(fmt.Errorf("unknown binding type: %#v", b))
	}

	// If no type was found, check if the type is allowed to be null
	if val == nil {
		for _, t := range types {
			if z, ok := t.(cwl.Null); ok {
				return bindings{
					{clb, z, nil, key, nil},
				}, nil
			}
		}
	}

	if val == nil {
		return nil, nil
	}

Loop:
	// an input descriptor describes multiple allowed types.
	// loop over the types, looking for the best match for the given input value.
	for _, t := range types {
		switch z := t.(type) {

		case cwl.InputArray:
			vals, ok := val.([]cwl.InputValue)
			if !ok {
				// input value is not an array.
				continue Loop
			}

			var out bindings

			for i, val := range vals {
				key := append(key, sortKey{getPos(z.InputBinding), i}...)
				b, err := e.walk(z, val, key)
				if err != nil {
					return nil, err
				}
				if b == nil {
					// array item values did not bind to the array descriptor.
					continue Loop
				}
				out = append(out, b...)
			}

			if out != nil {
				nested := make(bindings, len(out))
				copy(nested, out)
				b := &binding{clb, z, val, key, nested}
				// TODO revisit whether creating a nested tree (instead of flat) is always better/ok
				return bindings{b}, nil
			}

		case cwl.InputRecord:
			vals, ok := val.(map[string]cwl.InputValue)
			if !ok {
				// input value is not a record.
				continue Loop
			}

			var out bindings

			for i, field := range z.Fields {
				val, ok := vals[field.Name]
				// TODO lower case?
				if !ok {
					continue Loop
				}

				key := append(key, sortKey{getPos(field.InputBinding), i}...)
				b, err := e.walk(field, val, key)
				if err != nil {
					return nil, err
				}
				if b == nil {
					continue Loop
				}
				out = append(out, b...)
			}

			if out != nil {
				nested := make(bindings, len(out))
				copy(nested, out)
				b := &binding{clb, z, val, key, nested}
				out = append(out, b)
				return out, nil
			}

		case cwl.Boolean:
			if val == nil {
				continue Loop
			}
			v, err := cast.ToBoolE(val)
			if err != nil {
				continue Loop
			}
			return bindings{
				{clb, z, v, key, nil},
			}, nil

		case cwl.Int:
			v, err := cast.ToInt32E(val)
			if err != nil {
				continue Loop
			}
			return bindings{
				{clb, z, v, key, nil},
			}, nil

		case cwl.Long:
			v, err := cast.ToInt64E(val)
			if err != nil {
				continue Loop
			}
			return bindings{
				{clb, z, v, key, nil},
			}, nil

		case cwl.Float:
			v, err := cast.ToFloat32E(val)
			if err != nil {
				continue Loop
			}
			return bindings{
				{clb, z, v, key, nil},
			}, nil

		case cwl.Double:
			v, err := cast.ToFloat64E(val)
			if err != nil {
				continue Loop
			}
			return bindings{
				{clb, z, v, key, nil},
			}, nil

		case cwl.String:
			v, err := cast.ToStringE(val)
			if err != nil {
				continue Loop
			}

			return bindings{
				{clb, z, v, key, nil},
			}, nil

		case cwl.FileType:
			v, ok := val.(cwl.File)
			if !ok {
				continue Loop
			}
			f, err := ResolveFile(v, e.FS, clb.LoadContents)
			if err != nil {
				return nil, err
			}
			return bindings{
				{clb, z, *f, key, nil},
			}, nil

		case cwl.DirectoryType:
			v, ok := val.(cwl.Directory)
			if !ok {
				continue Loop
			}
			return bindings{
				{clb, z, v, key, nil},
			}, nil

		}
	}

	return nil, nil
}

func getPos(in *cwl.CommandLineBinding) int {
	if in == nil {
		return 0
	}
	return in.Position
}

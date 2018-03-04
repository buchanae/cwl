package cwl

import (
	"fmt"
	"github.com/buchanae/cwl/fs"
	"github.com/spf13/cast"
	"sort"
)

type Job struct {
	Command []string
	// TODO resource requests
	// TODO input files
	// TODO output binding description
	// environment
	// working directory
	// file staging
}

// TODO need to fully resolve file inputs (including secondary files)
//      before building job.

type Executor struct {
	FS fs.Filesystem
}

func NewExecutor() *Executor {
	return &Executor{FS: fs.NewLocal()}
}

// BuildJob builds command line arguments for an invocation a tool
// given a set of input values.
func (e *Executor) BuildJob(clt *CommandLineTool, vals InputValues) (*Job, error) {
	args := bindings{}

	// Add "arguments"
	for i, arg := range clt.Arguments {
		// TODO evaluate expressions
		b := &binding{arg, argType{}, string(arg.ValueFrom), sortKey{arg.Position, i}, nil}
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
func (e *Executor) walk(b bindable, val interface{}, key sortKey) (bindings, error) {
	types, clb := b.bindable()

	// If no type was found, check if the type is allowed to be null
	if val == nil {
		for _, t := range types {
			if z, ok := t.(Null); ok {
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

		case InputArray:
			vals, ok := val.([]InputValue)
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
				//out = append(out, b)
				//return out, nil
			}

		case InputRecord:
			vals, ok := val.(map[string]InputValue)
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

		case Boolean:
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

		case Int:
			v, err := cast.ToInt32E(val)
			if err != nil {
				continue Loop
			}
			return bindings{
				{clb, z, v, key, nil},
			}, nil

		case Long:
			v, err := cast.ToInt64E(val)
			if err != nil {
				continue Loop
			}
			return bindings{
				{clb, z, v, key, nil},
			}, nil

		case Float:
			v, err := cast.ToFloat32E(val)
			if err != nil {
				continue Loop
			}
			return bindings{
				{clb, z, v, key, nil},
			}, nil

		case Double:
			v, err := cast.ToFloat64E(val)
			if err != nil {
				continue Loop
			}
			return bindings{
				{clb, z, v, key, nil},
			}, nil

		case String:
			v, err := cast.ToStringE(val)
			if err != nil {
				continue Loop
			}

			return bindings{
				{clb, z, v, key, nil},
			}, nil

		case FileType:
			v, ok := val.(File)
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

		case DirectoryType:
			v, ok := val.(Directory)
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

func getPos(in *CommandLineBinding) int {
	if in == nil {
		return 0
	}
	return in.Position
}

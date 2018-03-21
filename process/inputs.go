package process

import (
	"github.com/buchanae/cwl"
	"github.com/spf13/cast"
)

/*** CWL input binding code ***/

// Binding binds an input type description (string, array, record, etc)
// to a concrete input value. this information is used while building
// command line args.
type Binding struct {
	clb *cwl.CommandLineBinding
	// the bound type (resolved by matching the input value to one of many allowed types)
	// can be nil, which means no matching type could be determined.
	Type interface{}
	// the value from the input object
	Value cwl.Value
	// used to determine the ordering of command line flags.
	// http://www.commonwl.org/v1.0/CommandLineTool.html#Input_binding
	sortKey sortKey
	nested  []*Binding
	name    string
}

// bindInput binds an input descriptor to a concrete value.
//
// bindInput is called recursively for types which have subtypes,
// such as array, record, etc.
//
// `name` is the field or parameter name.
// `types` is the list of types allowed by this input.
// `clb` is the cwl.CommandLineBinding describing how to bind this input.
// `val` is the input value for this input key.
// `key` is the sort key of the parent of this binding.
func (process *Process) bindInput(
	name string,
	types []cwl.InputType,
	clb *cwl.CommandLineBinding,
	secondaryFiles []cwl.Expression,
	val interface{},
	key sortKey,
) ([]*Binding, error) {

	// If no value was found, check if the type is allowed to be null.
	// If so, return a binding, otherwise fail.
	if val == nil {
		for _, t := range types {
			if z, ok := t.(cwl.Null); ok {
				return []*Binding{
					{clb, z, nil, key, nil, name},
				}, nil
			}
		}
		return nil, errf("failed to bind input, missing value")
	}

Loop:

	// An input descriptor describes multiple allowed types.
	// Loop over the types, looking for the best match for the given input value.
	for _, t := range types {
		switch z := t.(type) {

		case cwl.InputArray:
			vals, ok := val.([]cwl.Value)
			if !ok {
				// input value is not an array.
				continue Loop
			}

			// The input array is allowed to be empty,
			// so this must be a non-nil slice.
			out := []*Binding{}

			for i, val := range vals {
				subkey := append(key, sortKey{getPos(z.InputBinding), i}...)
				b, err := process.bindInput("", z.Items, z.InputBinding, nil, val, subkey)
				if err != nil {
					return nil, err
				}
				if b == nil {
					// array item values did not bind to the array descriptor.
					continue Loop
				}
				out = append(out, b...)
			}

			nested := make([]*Binding, len(out))
			copy(nested, out)
			b := &Binding{clb, z, val, key, nested, name}
			// TODO revisit whether creating a nested tree (instead of flat) is always better/ok
			return []*Binding{b}, nil

		case cwl.InputRecord:
			vals, ok := val.(map[string]cwl.Value)
			if !ok {
				// input value is not a record.
				continue Loop
			}

			var out []*Binding

			for i, field := range z.Fields {
				val, ok := vals[field.Name]
				// TODO lower case?
				if !ok {
					continue Loop
				}

				subkey := append(key, sortKey{getPos(field.InputBinding), i}...)
				b, err := process.bindInput(field.Name, field.Type, field.InputBinding, nil, val, subkey)
				if err != nil {
					return nil, err
				}
				if b == nil {
					continue Loop
				}
				out = append(out, b...)
			}

			if out != nil {
				nested := make([]*Binding, len(out))
				copy(nested, out)
				b := &Binding{clb, z, val, key, nested, name}
				out = append(out, b)
				return out, nil
			}

		case cwl.Any:
			return []*Binding{
				{clb, z, val, key, nil, name},
			}, nil

		case cwl.Boolean:
			v, err := cast.ToBoolE(val)
			if err != nil {
				continue Loop
			}
			return []*Binding{
				{clb, z, v, key, nil, name},
			}, nil

		case cwl.Int:
			v, err := cast.ToInt32E(val)
			if err != nil {
				continue Loop
			}
			return []*Binding{
				{clb, z, v, key, nil, name},
			}, nil

		case cwl.Long:
			v, err := cast.ToInt64E(val)
			if err != nil {
				continue Loop
			}
			return []*Binding{
				{clb, z, v, key, nil, name},
			}, nil

		case cwl.Float:
			v, err := cast.ToFloat32E(val)
			if err != nil {
				continue Loop
			}
			return []*Binding{
				{clb, z, v, key, nil, name},
			}, nil

		case cwl.Double:
			v, err := cast.ToFloat64E(val)
			if err != nil {
				continue Loop
			}
			return []*Binding{
				{clb, z, v, key, nil, name},
			}, nil

		case cwl.String:
			v, err := cast.ToStringE(val)
			if err != nil {
				continue Loop
			}

			return []*Binding{
				{clb, z, v, key, nil, name},
			}, nil

		case cwl.FileType:
			v, ok := val.(cwl.File)
			if !ok {
				continue Loop
			}

			f, err := process.resolveFile(v, clb.GetLoadContents())
			if err != nil {
				return nil, err
			}
			// TODO figure out a good way to do this.
			f.Path = "/inputs/" + f.Path
			for _, expr := range secondaryFiles {
				process.resolveSecondaryFiles(f, expr)
			}

			return []*Binding{
				{clb, z, f, key, nil, name},
			}, nil

		case cwl.DirectoryType:
			v, ok := val.(cwl.Directory)
			if !ok {
				continue Loop
			}
			// TODO resolve directory
			return []*Binding{
				{clb, z, v, key, nil, name},
			}, nil

		}
	}

	return nil, errf("failed to bind input, missing value")
}

package cwllib

import (
	"github.com/buchanae/cwl"
	"github.com/spf13/cast"
)

/*** CWL input binding code ***/

// binding binds an input type description (string, array, record, etc)
// to a concrete input value. this information is used while building
// command line args.
type binding struct {
	clb *cwl.CommandLineBinding
	// the bound type (resolved by matching the input value to one of many allowed types)
	// can be nil, which means no matching type could be determined.
	typ interface{}
	// the value from the input object
	value cwl.Value
	// used to determine the ordering of command line flags.
	// http://www.commonwl.org/v1.0/CommandLineTool.html#Input_binding
	sortKey sortKey
	nested  []*binding
}

// bindInput binds an input descriptor to a concrete value.
//
// bindInput is called recursively for types which have subtypes,
// such as array, record, etc.
//
// `fs` provides access to the filesystem.
// `types` is the list of types allowed by this input.
// `clb` is the cwl.CommandLineBinding describing how to bind this input.
// `val` is the input value for this input key.
// `key` is the sort key of the parent of this binding.
func (job *Job) bindInput(
  types []cwl.InputType,
  clb *cwl.CommandLineBinding,
  secondaryFiles []cwl.Expression,
  val interface{},
  key sortKey,
) ([]*binding, error) {

	// If no value was found, check if the type is allowed to be null.
  // If so, return a binding.
	if val == nil {
		for _, t := range types {
			if z, ok := t.(cwl.Null); ok {
				return []*binding{
					{clb, z, nil, key, nil},
				}, nil
			}
		}
	}

  // TODO maybe return an error here to be more explicit.
	if val == nil {
		return nil, nil
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

			var out []*binding

			for i, val := range vals {
				key := append(key, sortKey{getPos(z.InputBinding), i}...)
				b, err := job.bindInput(z.Items, z.InputBinding, nil, val, key)
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
				nested := make([]*binding, len(out))
				copy(nested, out)
				b := &binding{clb, z, val, key, nested}
				// TODO revisit whether creating a nested tree (instead of flat) is always better/ok
				return []*binding{b}, nil
			}

		case cwl.InputRecord:
			vals, ok := val.(map[string]cwl.Value)
			if !ok {
				// input value is not a record.
				continue Loop
			}

			var out []*binding

			for i, field := range z.Fields {
				val, ok := vals[field.Name]
				// TODO lower case?
				if !ok {
					continue Loop
				}

				key := append(key, sortKey{getPos(field.InputBinding), i}...)
				b, err := job.bindInput(field.Type, field.InputBinding, nil, val, key)
				if err != nil {
					return nil, err
				}
				if b == nil {
					continue Loop
				}
				out = append(out, b...)
			}

			if out != nil {
				nested := make([]*binding, len(out))
				copy(nested, out)
				b := &binding{clb, z, val, key, nested}
				out = append(out, b)
				return out, nil
			}

		case cwl.Boolean:
      // TODO if-statement above means val should never be nil at this point?
			if val == nil {
				continue Loop
			}
			v, err := cast.ToBoolE(val)
			if err != nil {
				continue Loop
			}
			return []*binding{
				{clb, z, v, key, nil},
			}, nil

		case cwl.Int:
			v, err := cast.ToInt32E(val)
			if err != nil {
				continue Loop
			}
			return []*binding{
				{clb, z, v, key, nil},
			}, nil

		case cwl.Long:
			v, err := cast.ToInt64E(val)
			if err != nil {
				continue Loop
			}
			return []*binding{
				{clb, z, v, key, nil},
			}, nil

		case cwl.Float:
			v, err := cast.ToFloat32E(val)
			if err != nil {
				continue Loop
			}
			return []*binding{
				{clb, z, v, key, nil},
			}, nil

		case cwl.Double:
			v, err := cast.ToFloat64E(val)
			if err != nil {
				continue Loop
			}
			return []*binding{
				{clb, z, v, key, nil},
			}, nil

		case cwl.String:
			v, err := cast.ToStringE(val)
			if err != nil {
				continue Loop
			}

			return []*binding{
				{clb, z, v, key, nil},
			}, nil

		case cwl.FileType:
			v, ok := val.(cwl.File)
			if !ok {
				continue Loop
			}

			f, err := job.resolveFile(v, clb.LoadContents)
			if err != nil {
				return nil, err
			}
      for _, expr := range secondaryFiles {
        job.resolveSecondaryFiles(f, expr)
      }

			return []*binding{
				{clb, z, *f, key, nil},
			}, nil

		case cwl.DirectoryType:
			v, ok := val.(cwl.Directory)
			if !ok {
				continue Loop
			}
      // TODO resolve directory
			return []*binding{
				{clb, z, v, key, nil},
			}, nil

		}
	}

	return nil, nil
}

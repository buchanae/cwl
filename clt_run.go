package cwl

import (
	"fmt"
	"github.com/spf13/cast"
	"sort"
)

func buildCommand(clt *CommandLineTool, vals map[string]interface{}) ([]string, error) {
	arr := bindings{}

	for i, arg := range clt.Arguments {
		// TODO evaluate expressions
		b := &binding{arg, argType{}, string(arg.ValueFrom), sortKey{arg.Position, i}}
		arr = append(arr, b)
	}

	for i, in := range clt.Inputs {
		k := sortKey{in.InputBinding.Position, i}

		b := walk(in, vals[in.ID], k)
		if b == nil {
			return nil, fmt.Errorf("no valid binding found for input: %s, %s", in.ID)
		}
		arr = append(arr, b...)
	}

	sort.Stable(arr)
	debug(arr)

	// Now collect the input bindings into command line arguments
	args := append([]string{}, clt.BaseCommand...)
	for _, b := range arr {
		args = append(args, b.args()...)
	}

	return args, nil
}

func walk(b bindable, val interface{}, key sortKey) bindings {
	types, clb := b.bindable()

Loop:
	for _, t := range types {
		switch z := t.(type) {

		case InputArray:
			vals, ok := val.([]interface{})
			if !ok {
				continue Loop
			}

			var out bindings

			for i, val := range vals {
				key := append(key, sortKey{z.InputBinding.Position, i}...)
				b := walk(z, val, key)
				if b == nil {
					continue Loop
				}
				out = append(out, b...)
			}
			if out != nil {
				out = append(out, &binding{clb, z, val, key})
				return out
			}

		case InputRecord:
			vals, ok := val.(map[string]interface{})
			if !ok {
				continue Loop
			}

			var out bindings

			for i, field := range z.Fields {
				val, ok := vals[field.Name]
				// TODO lower case?
				if !ok {
					continue Loop
				}

				key := append(key, sortKey{field.InputBinding.Position, i}...)
				b := walk(field, val, key)
				if b == nil {
					continue Loop
				}
				out = append(out, b...)
			}
			if out != nil {
				return out
			}

		case Boolean:
			if val == nil {
				continue Loop
			}
			_, err := cast.ToBoolE(val)
			if err != nil {
				continue Loop
			}
			return bindings{
				{clb, z, val, key},
			}

		case Int:
			_, err := cast.ToInt32E(val)
			if err != nil {
				continue Loop
			}
			return bindings{
				{clb, z, val, key},
			}

		case Long:
			_, err := cast.ToInt64E(val)
			if err != nil {
				continue Loop
			}
			return bindings{
				{clb, z, val, key},
			}

		case Float:
			_, err := cast.ToFloat32E(val)
			if err != nil {
				continue Loop
			}
			return bindings{
				{clb, z, val, key},
			}

		case Double:
			_, err := cast.ToFloat64E(val)
			if err != nil {
				continue Loop
			}
			return bindings{
				{clb, z, val, key},
			}

		case String:
			if val == nil {
				continue Loop
			}

			return bindings{
				{clb, z, val, key},
			}

		case FileType:
			// TODO need to get map and unmarshal (loader?)
			//      into File struct

		case DirectoryType:
			// TODO need to get map and unmarshal (loader?)
			//      into Directory struct

		}
	}

	// If no type was found, check if the type is allowed to be null
	if val == nil {
		for _, t := range types {
			if z, ok := t.(Null); ok {
				return bindings{
					{clb, z, nil, key},
				}
			}
		}
	}

	return nil
}

func prefixArg(prefix, arg string, sep bool) []string {
	if prefix == "" {
		return []string{arg}
	}
	if sep {
		return []string{prefix, arg}
	}
	return []string{prefix + arg}
}

package cwl

import (
	"fmt"
	"github.com/spf13/cast"
	"sort"
)

func buildCommand(clt CommandLineTool, vals map[string]interface{}) error {
	arr := bindings{}

	// collect the tree of CommandLineBindings
	for i, arg := range clt.Arguments {
		arr = append(arr, &binding{clb: arg, idx: i})
	}

	// bind inputs parameter definitions to the input object.
	// this binds the input to one of potentially many types.
	for i, in := range clt.Inputs {
		b, err := bindInput(in, vals)
		if err != nil {
			return err
		}
		if b == nil {
			// binding is allowed to be "null" so skip it
			continue
		}
		b.idx = i
		arr = append(arr, b)
	}

	sort.Stable(arr)

	// Now collect the input bindings into command line arguments
	args := append([]string{}, clt.BaseCommand...)
	for _, arg := range arr {
		args = append(args, arg.args...)
	}

	fmt.Println(args)
	return nil
}

func bindInput(in CommandInput, vals map[string]interface{}) (*binding, error) {

	val, ok := vals[in.ID]
	if !ok {
		// the input parameter is missing a matching concrete value.
		// check to see if "null" is in the parameter type list.
		for _, ty := range in.Type {
			if _, ok := ty.(Null); ok {
				// it's ok that this field is null.
				return nil, nil
			}
		}
		return nil, fmt.Errorf("missing input value")
	}

	b := &binding{clb: in.InputBinding}
	strval := fmt.Sprintf("%s", val)

	// the input parameter can have multiple allowed types.
	// try to find one that matches the concrete value.
	for _, ty := range in.Type {

		// handle complex types first
		switch ty.(type) {
		case FileType:
			// TODO need to get map and unmarshal (loader?)
			//      into File struct

		case DirectoryType:
			// TODO need to get map and unmarshal (loader?)
			//      into Directory struct

		case ArrayType:
			arrval, ok := val.([]interface{})
			if !ok {
				continue
			}
			_ = arrval

		case RecordType:
			mapval, ok := val.(map[string]interface{})
			if !ok {
				continue
			}
			_ = mapval
		}

		// now handle primitive types
		switch ty.(type) {
		case Boolean:
			bv, err := cast.ToBoolE(val)
			if err != nil {
				continue
			}
			if !bv {
				return nil, nil
			}
			if b.clb.Prefix == "" {
				return nil, fmt.Errorf("boolean value without prefix")
			}
			b.args = []string{b.clb.Prefix}
			return b, nil

		case Int:
			_, err := cast.ToInt32E(val)
			if err != nil {
				continue
			}
			b.args = prefixArg(b.clb.Prefix, strval, b.clb.Separate)
			return b, nil

		case Long:
			_, err := cast.ToInt64E(val)
			if err != nil {
				continue
			}
			b.args = prefixArg(b.clb.Prefix, strval, b.clb.Separate)
			return b, nil

		case Float:
			_, err := cast.ToFloat32E(val)
			if err != nil {
				continue
			}
			b.args = prefixArg(b.clb.Prefix, strval, b.clb.Separate)
			return b, nil

		case Double:
			_, err := cast.ToFloat64E(val)
			if err != nil {
				continue
			}
			b.args = prefixArg(b.clb.Prefix, strval, b.clb.Separate)
			return b, nil

		case String:
			b.args = prefixArg(b.clb.Prefix, strval, b.clb.Separate)
			return b, nil
		}
	}

	return nil, fmt.Errorf("no matching type found")
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

// binding collects information used while building
// a concrete command line.
type binding struct {
	clb CommandLineBinding
	// index of the binding in an array
	idx  int
	key  string
	args []string
}

// bindings defines the rules for sorting bindings;
// http://www.commonwl.org/v1.0/CommandLineTool.html#Input_binding
type bindings []*binding

func (s bindings) Len() int      { return len(s) }
func (s bindings) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s bindings) Less(i, j int) bool {
	return s[i].clb.Position < s[j].clb.Position || s[i].idx < s[j].idx
}

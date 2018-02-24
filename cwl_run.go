package cwl

import (
	"fmt"
	"github.com/spf13/cast"
	"sort"
)

func buildCommand(clt CommandLineTool, vals map[string]interface{}) error {
	arr := bindings{}

	for i, arg := range clt.Arguments {
		b := &binding{
			clb:     arg,
			typ:     argType{},
			sortKey: sortKey{arg.Position, i},
		}
		arr = append(arr, flatten(b)...)
	}

	for i, in := range clt.Inputs {
		b := &binding{
			clb:     in.InputBinding,
			value:   vals[in.ID],
			sortKey: sortKey{in.InputBinding.Position, i},
			typ:     matchType(in.Type, vals[in.ID]),
		}
		arr = append(arr, b)
	}

	sort.Stable(arr)

	// Now collect the input bindings into command line arguments
	args := append([]string{}, clt.BaseCommand...)

	fmt.Println(args)
	return nil
}

// flatten flattens nested array and record types into a flat list of bindings.
func flatten(b *binding) []*binding {
	arr := []*binding{b}

	switch t := b.typ.(type) {
	case InputArray:

		vals := b.value.([]interface{})
		for i, val := range vals {
			a := &binding{
				clb:     t.InputBinding,
				value:   val,
				sortKey: append(b.sortKey, sortKey{t.InputBinding.Position, i}...),
				typ:     matchType(t.Items, val),
			}
			arr = append(arr, flatten(a)...)
		}

	case InputRecord:
	}
	return arr
}

// binding binds an input type description (string, array, record, etc)
// to a concrete input value. this information is used while building
// command line args.
type binding struct {
	clb CommandLineBinding
	// the bound type (resolved by matching the input value to one of many allowed types)
	// can be nil, which means no matching type could be determined.
	typ InputType
	// the value from the input object
	value interface{}
	// used to determine the ordering of command line flags.
	// http://www.commonwl.org/v1.0/CommandLineTool.html#Input_binding
	sortKey sortKey
}

// matchType matches the input value to one of possibly many types
// allowed by the input parameter specification.
// returns nil if no matching type is found.
func matchType(types []InputType, val interface{}) InputType {
	for _, typ := range types {

		// handle complex types first
		switch typ.(type) {
		case FileType:
			// TODO need to get map and unmarshal (loader?)
			//      into File struct

		case DirectoryType:
			// TODO need to get map and unmarshal (loader?)
			//      into Directory struct

		case InputArray:
			_, ok := val.([]interface{})
			if !ok {
				continue
			}
			return typ

		case InputRecord:
			_, ok := val.(map[string]interface{})
			if !ok {
				continue
			}
			return typ

		case Boolean:
			_, err := cast.ToBoolE(val)
			if err != nil {
				continue
			}
			return typ
			/*
				if !bv {
					return nil, nil
				}
				if b.clb.Prefix == "" {
					return nil, fmt.Errorf("boolean value without prefix")
				}
				b.args = []string{b.clb.Prefix}
				return b, nil
			*/

		case Int:
			_, err := cast.ToInt32E(val)
			if err != nil {
				continue
			}
			return typ
			//b.args = prefixArg(b.clb.Prefix, strval, b.clb.Separate)

		case Long:
			_, err := cast.ToInt64E(val)
			if err != nil {
				continue
			}
			return typ
			//b.args = prefixArg(b.clb.Prefix, strval, b.clb.Separate)
			//return b, nil

		case Float:
			_, err := cast.ToFloat32E(val)
			if err != nil {
				continue
			}
			return typ
			//b.args = prefixArg(b.clb.Prefix, strval, b.clb.Separate)
			//return b, nil

		case Double:
			_, err := cast.ToFloat64E(val)
			if err != nil {
				continue
			}
			return typ
			//b.args = prefixArg(b.clb.Prefix, strval, b.clb.Separate)

		case String:
			return typ
			//b.args = prefixArg(b.clb.Prefix, strval, b.clb.Separate)
			//return b, nil

		case Null:
			if val != nil {
				continue
			}
			return typ
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

type sortKey []interface{}

// bindings defines the rules for sorting bindings;
// http://www.commonwl.org/v1.0/CommandLineTool.html#Input_binding
type bindings []*binding

func (s bindings) Len() int      { return len(s) }
func (s bindings) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s bindings) Less(i, j int) bool {
	z := compareKey(s[i].sortKey, s[j].sortKey)
	return z == -1
}

// compare two sort keys.
//
// The result will be 0 if i==j, -1 if i < j, and +1 if i > j.
func compareKey(i, j sortKey) int {
	for x := 0; x < len(i) || x < len(j); x++ {
		if x >= len(i) {
			// i key is shorter than j
			return -1
		}
		if x >= len(j) {
			// j key is shorter than i
			return 1
		}
		z := compare(i[x], j[x])
		if z != 0 {
			return z
		}
	}
	return 0
}

// compare two sort key items, because sort keys may have mixed ints and strings.
// cwl spec: "ints sort before strings", i.e all ints are less than all strings.
//
// The result will be 0 if i==j, -1 if i < j, and +1 if i > j.
func compare(iv, jv interface{}) int {
	istr, istrok := iv.(string)
	jstr, jstrok := jv.(string)
	iint, iintok := iv.(int)
	jint, jintok := jv.(int)

	switch {
	case istrok && jintok:
		// i is a string, j is an int
		// cwl spec: "ints sort before strings"
		return 1
	case iintok && jstrok:
		// i is an int, j is a string
		// cwl spec: "ints sort before strings"
		return -1

	// both are strings
	case istrok && jstrok && istr == jstr:
		return 0
	case istrok && jstrok && istr < jstr:
		return -1
	case istrok && jstrok && istr > jstr:
		return 1

	// both are ints
	case iintok && jintok && iint == jint:
		return 0
	case iintok && jintok && iint < jint:
		return -1
	case iintok && jintok && iint > jint:
		return 1
	}
	return 0
}

// argType is used internally to mark a binding as coming from "CommandLineTool.Arguments"
type argType struct{}

func (argType) inputtype()     {}
func (argType) String() string { return "argument" }

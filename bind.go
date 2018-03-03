package cwl

import (
	"fmt"
	"strings"
)

// bindable describes a type which may have input values bound to it:
// CommandInput, InputArray, InputField
type bindable interface {
	bindable() ([]InputType, *CommandLineBinding)
}

// binding binds an input type description (string, array, record, etc)
// to a concrete input value. this information is used while building
// command line args.
type binding struct {
	clb *CommandLineBinding
	// the bound type (resolved by matching the input value to one of many allowed types)
	// can be nil, which means no matching type could be determined.
	typ InputType
	// the value from the input object
	value interface{}
	// used to determine the ordering of command line flags.
	// http://www.commonwl.org/v1.0/CommandLineTool.html#Input_binding
	sortKey sortKey
	nested  bindings
}

// args converts a binding into a list of formatted command line arguments.
func (b *binding) args() []string {
	switch b.typ.(type) {

	case InputArray:
		strval := ""
		if b.clb != nil && b.clb.ItemSeparator != "" {
			var nested []string
			for _, nb := range b.nested {
				nested = append(nested, nb.args()...)
			}
			strval = strings.Join(nested, b.clb.ItemSeparator)
		}
		return formatArg(b.clb, strval)

	case InputRecord:
		return formatArg(b.clb, "")

	case argType:
		strval := fmt.Sprintf("%s", b.value)
		return formatArg(b.clb, strval)

	case String, Int, Long, Float, Double:
		strval := fmt.Sprintf("%s", b.value)
		return formatArg(b.clb, strval)

	case FileType:
		f := b.value.(File)
		return formatArg(b.clb, f.Path)

	case DirectoryType:
		d := b.value.(File)
		return formatArg(b.clb, d.Path)

	case Boolean:
		bv := b.value.(bool)
		if !bv {
			return nil
		}
		/*
		   TODO find a place for this validation
		   if b.clb.Prefix == "" {
		     return nil, fmt.Errorf("boolean value without prefix")
		   }
		*/
		prefix := ""
		if b.clb != nil {
			prefix = b.clb.Prefix
		}
		return []string{prefix}
	}
	return nil
}

// formatArg applies some command line binding rules to a CLI argument,
// such as flag prefixing, separator, etc.
// http://www.commonwl.org/v1.0/CommandLineTool.html#CommandLineBinding
func formatArg(clb *CommandLineBinding, arg string) []string {
	sep := true
	prefix := ""

	if clb != nil {
		prefix = clb.Prefix
		sep = clb.Separate()
	}

	switch {
	case prefix == "" && arg == "":
		return nil

	case prefix == "":
		return []string{arg}

	case arg == "":
		return []string{prefix}

	case sep:
		return []string{prefix, arg}
	default:
		return []string{prefix + arg}
	}
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

func (c CommandInput) bindable() ([]InputType, *CommandLineBinding) {
	return c.Type, c.InputBinding
}
func (i InputArray) bindable() ([]InputType, *CommandLineBinding) {
	return i.Items, i.InputBinding
}
func (i InputField) bindable() ([]InputType, *CommandLineBinding) {
	return i.Type, i.InputBinding
}
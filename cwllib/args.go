package cwllib

import (
	"fmt"
	"github.com/buchanae/cwl"
	"sort"
	"strings"
)

/*** CWL tool command line argument building code ***/

func (job *Job) Command() ([]string, error) {

	args := make([]*binding, len(job.bindings))
	copy(args, job.bindings)

	// Add "CommandLineTool.arguments"
	for i, arg := range job.tool.Arguments {
		if arg.ValueFrom == "" {
			return nil, errf("valueFrom is required but missing for argument %d", i)
		}
		args = append(args, &binding{
			arg, argType{}, nil, sortKey{arg.Position, i}, nil,
		})
	}

	// Evaluate "valueFrom" expression.
	for _, b := range args {
		if b.clb.ValueFrom != "" {
			val, err := job.eval(b.clb.ValueFrom, b.value)
			if err != nil {
				return nil, errf("failed to eval argument value: %s", err)
			}
			b.value = val
		}
	}

	sort.Stable(bySortKey(args))

	// Now collect the input bindings into command line arguments
	cmd := append([]string{}, job.tool.BaseCommand...)
	for _, b := range args {
		cmd = append(cmd, bindArgs(b)...)
	}

	return cmd, nil
}

// args converts a binding into a list of formatted command line arguments.
func bindArgs(b *binding) []string {
	switch b.typ.(type) {

	case cwl.InputArray:
		// cwl spec:
		// "If itemSeparator is specified, add prefix and the join the array
		// into a single string with itemSeparator separating the items..."
		if b.clb != nil && b.clb.ItemSeparator != "" {

			var nested []cwl.Value
			for _, nb := range b.nested {
				nested = append(nested, nb.value)
			}
			return formatArgs(b.clb, nested...)

			// cwl spec:
			// "...otherwise first add prefix, then recursively process individual elements."
		} else {
			args := formatArgs(b.clb)

			for _, nb := range b.nested {
				args = append(args, bindArgs(nb)...)
			}
			return args
		}

	case cwl.InputRecord:
		// TODO

	case cwl.String, cwl.Int, cwl.Long, cwl.Float, cwl.Double, cwl.FileType,
		cwl.DirectoryType, argType:
		return formatArgs(b.clb, b.value)

	case cwl.Boolean:
		// cwl spec:
		// "boolean: If true, add prefix to the command line. If false, add nothing."
		bv := b.value.(bool)
		if bv && b.clb != nil && b.clb.Prefix != "" {
			return formatArgs(b.clb)
		}
	}
	return nil
}

// formatArgs applies some command line binding rules to a CLI argument,
// such as prefix, separate, etc.
// http://www.commonwl.org/v1.0/CommandLineTool.html#CommandLineBinding
func formatArgs(clb *cwl.CommandLineBinding, args ...cwl.Value) []string {
	sep := true
	prefix := ""
	join := ""

	if clb != nil {
		prefix = clb.Prefix
		sep = clb.Separate.Value()
		join = clb.ItemSeparator
	}

	var strargs []string
	for _, arg := range args {
		strargs = append(strargs, valueToStrings(arg)...)
	}

	if join != "" && strargs != nil {
		strargs = []string{strings.Join(strargs, join)}
	}

	if prefix != "" {
		if sep {
			strargs = append([]string{prefix}, strargs...)
		} else if strargs != nil {
			strargs[0] = prefix + strargs[0]
		} else {
			strargs = []string{prefix}
		}
	}
	return strargs
}

func valueToStrings(v cwl.Value) []string {
	switch z := v.(type) {
	case []interface{}:
		var out []string
		for _, v := range z {
			out = append(out, valueToStrings(v)...)
		}
		return out
	case int, int32, int64, float32, float64, bool, string:
		return []string{fmt.Sprintf("%v", z)}
	case cwl.File:
		return []string{z.Path}
	case cwl.Directory:
		return []string{z.Path}
	}
	return nil
}

type sortKey []interface{}

// bySortKey defines the rules for sorting bindings;
// http://www.commonwl.org/v1.0/CommandLineTool.html#Input_binding
type bySortKey []*binding

func (s bySortKey) Len() int      { return len(s) }
func (s bySortKey) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s bySortKey) Less(i, j int) bool {
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

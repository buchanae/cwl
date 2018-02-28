package cwl

import (
	"fmt"
	"github.com/commondream/yamlast"
	"strings"
)

func (l *loader) MappingToDocument(n node) (Document, error) {

	class := findKey(n, "class")
	switch class {

	case "commandlinetool":
		t := &CommandLineTool{}
		if err := l.load(n, t); err != nil {
			return nil, err
		}
		return t, nil

	case "workflow":
		wf := &Workflow{}
		if err := l.load(n, wf); err != nil {
			return nil, err
		}
		return wf, nil

	default:
		return nil, fmt.Errorf("unknown document class: '%s'", class)
	}
}

func (l *loader) ScalarToDocument(n node) (Document, error) {
	return DocumentRef{URL: n.Value}, nil
}

func (l *loader) SeqToAny(n node) (Any, error) {
	return nil, nil
}

func (l *loader) SeqToStringSlice(n node) ([]string, error) {
	strs := []string{}
	for _, c := range n.Children {
		strs = append(strs, c.Value)
	}
	return strs, nil
}

func (l *loader) SeqToCommandLineBindingSlice(n node) ([]CommandLineBinding, error) {
	var b []CommandLineBinding
	for _, c := range n.Children {
		if c.Kind != yamlast.MappingNode {
			return nil, fmt.Errorf("unhandled command line binding type")
		}
		clb := CommandLineBinding{}
		err := l.load(c, &clb)
		if err != nil {
			return nil, err
		}
		b = append(b, clb)
	}
	return b, nil
}

func (l *loader) SeqToString(n node) (string, error) {
	s := ""
	for _, c := range n.Children {
		if c.Kind != yamlast.ScalarNode {
			return "", fmt.Errorf("unhandled string concat type")
		}
		if s != "" {
			s += "\n" + c.Value
		} else {
			s = c.Value
		}
	}
	return s, nil
}

func (l *loader) MappingToInputTypeSlice(n node) ([]InputType, error) {
	t, err := l.MappingToInputType(n)
	if err != nil {
		return nil, err
	}
	return []InputType{t}, nil
}

func (l *loader) MappingToOutputTypeSlice(n node) ([]OutputType, error) {
	t, err := l.MappingToOutputType(n)
	if err != nil {
		return nil, err
	}
	return []OutputType{t}, nil
}

func (l *loader) SeqToInputTypeSlice(n node) ([]InputType, error) {
	var out []InputType
	for _, c := range n.Children {
		var t InputType
		err := l.load(c, &t)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

func (l *loader) MappingToInputType(n node) (InputType, error) {
	typ := findKey(n, "type")
	switch typ {
	case "array":
		a := InputArray{}
		err := l.load(n, &a)
		return a, err
	case "record":
		rec := InputRecord{}
		err := l.load(n, &rec)
		return rec, err
	case "enum":
		return InputEnum{}, nil
	}
	panic("unknown type")
}

func (l *loader) MappingToOutputType(n node) (OutputType, error) {
	typ := findKey(n, "type")
	switch typ {
	case "array":
		a := OutputArray{}
		err := l.load(n, &a)
		return a, err
	case "record":
		rec := OutputRecord{}
		err := l.load(n, &rec)
		return rec, err
	case "enum":
		return OutputEnum{}, nil
	}
	panic("unknown type")
}

func (l *loader) MappingToInputFieldSlice(n node) ([]InputField, error) {
	var fields []InputField

	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v
		i := InputField{Name: k}
		err := l.load(v, &i)
		if err != nil {
			return nil, err
		}
		fields = append(fields, i)
	}
	return fields, nil
}

func (l *loader) ScalarToInputTypeSlice(n node) ([]InputType, error) {

	if strings.HasSuffix(n.Value, "?") {
		name := strings.TrimSuffix(n.Value, "?")
		t, ok := getInputTypeByName(name)
		if ok {
			return []InputType{t, Null{}}, nil
		}
	}

	t, err := l.ScalarToInputType(n)
	if err != nil {
		return nil, err
	}
	if t != nil {
		return []InputType{t}, nil
	}
	panic("unhandled cwl type")
}

func (l *loader) ScalarToOutputTypeSlice(n node) ([]OutputType, error) {

	if strings.HasSuffix(n.Value, "?") {
		name := strings.TrimSuffix(n.Value, "?")
		t, ok := getOutputTypeByName(name)
		if ok {
			return []OutputType{t, Null{}}, nil
		}
	}

	t, err := l.ScalarToOutputType(n)
	if err != nil {
		return nil, err
	}
	if t != nil {
		return []OutputType{t}, nil
	}
	panic("unhandled cwl type")
}

// TODO is "string[]?" acceptable?
func (l *loader) ScalarToInputType(n node) (InputType, error) {
	name := n.Value

	if strings.HasSuffix(name, "[]") {
		name := strings.TrimSuffix(name, "[]")
		t, ok := getInputTypeByName(name)
		if ok {
			return InputArray{Items: []InputType{t}}, nil
		}
	}

	t, ok := getInputTypeByName(name)
	if ok {
		return t, nil
	}
	// TODO should convert an unknown node into an IRI type reference
	return nil, fmt.Errorf("unhandled scalar type: %s", n.Value)
}

// TODO is "string[]?" acceptable?
func (l *loader) ScalarToOutputType(n node) (OutputType, error) {
	name := n.Value

	if strings.HasSuffix(name, "[]") {
		name := strings.TrimSuffix(name, "[]")
		t, ok := getOutputTypeByName(name)
		if ok {
			return OutputArray{Items: []OutputType{t}}, nil
		}
	}

	t, ok := getOutputTypeByName(name)
	if ok {
		return t, nil
	}
	// TODO should convert an unknown node into an IRI type reference
	return nil, fmt.Errorf("unhandled scalar type: %s", n.Value)
}

func (l *loader) SeqToExpressionSlice(n node) ([]Expression, error) {
	exprs := []Expression{}
	for _, c := range n.Children {
		e := Expression("")
		l.load(c, &e)
		exprs = append(exprs, e)
	}
	return exprs, nil
}

func (l *loader) ScalarToExpressionSlice(n node) ([]Expression, error) {
	return []Expression{Expression(n.Value)}, nil
}

var inputTypesByName = map[string]InputType{}

func init() {
	ts := []InputType{
		Null{}, Boolean{}, Int{}, Long{}, Float{}, Double{}, String{},
		FileType{}, DirectoryType{}, InputRecord{}, InputArray{}, InputEnum{},
	}
	for _, t := range ts {
		name := strings.ToLower(t.String())
		inputTypesByName[name] = t
	}
}

func getInputTypeByName(name string) (InputType, bool) {
	name = strings.ToLower(name)
	t, ok := inputTypesByName[name]
	return t, ok
}

var outputTypesByName = map[string]OutputType{}

func init() {
	ts := []OutputType{
		Null{}, Boolean{}, Int{}, Long{}, Float{}, Double{}, String{},
		FileType{}, DirectoryType{}, OutputRecord{}, OutputArray{}, OutputEnum{},
	}
	for _, t := range ts {
		name := strings.ToLower(t.String())
		outputTypesByName[name] = t
	}
}

func getOutputTypeByName(name string) (OutputType, bool) {
	name = strings.ToLower(name)
	t, ok := outputTypesByName[name]
	return t, ok
}

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
	b := []CommandLineBinding{}
	for _, c := range n.Children {
		if c.Kind != yamlast.MappingNode {
			return nil, fmt.Errorf("unhandled command line binding type")
		}
		clb := CommandLineBinding{
			Separate: true,
		}
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

func (l *loader) MappingToTypeSlice(n node) ([]Type, error) {
	t, err := l.MappingToType(n)
	if err != nil {
		return nil, err
	}
	return []Type{t.(Type)}, nil
}

func (l *loader) MappingToType(n node) (Type, error) {
	typ := findKey(n, "type")
	switch typ {
	case "array":
		i, ok := findValue(n, "items")
		if !ok {
			panic("")
		}
		a := InputArray{}
		err := l.load(i, &a.Items)
		return a, err
	case "record":
		return InputRecord{}, nil
	case "enum":
		return InputEnum{}, nil
	}
	panic("unknown type")
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

func getInputTypeByName(name string) (InputType, bool) {
	name = strings.ToLower(name)
	t, ok := inputTypesByName[name]
	return t, ok
}

func getOutputTypeByName(name string) (OutputType, bool) {
	name = strings.ToLower(name)
	t, ok := outputTypesByName[name]
	return t, ok
}

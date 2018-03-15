package cwl

import (
	"fmt"
	"regexp"
)

// YAML allows all these values for booleans
// ...because the world needed more than true/false?
var yamlTrueRX = regexp.MustCompile(`y|Y|yes|Yes|YES|true|True|TRUE|on|On|ON`)
var yamlFalseRX = regexp.MustCompile(`n|N|no|No|NO|false|False|FALSE|off|Off|OFF`)

func (l *loader) ScalarToCommandOutput(n node) (CommandOutput, error) {
	o := CommandOutput{}
	err := l.load(n, &o.Type)
	return o, err
}

func (l *loader) ScalarToCommandInput(n node) (CommandInput, error) {
	o := CommandInput{}
	err := l.load(n, &o.Type)
	return o, err
}

func (l *loader) ScalarToCommandLineBinding(n node) (CommandLineBinding, error) {
	return CommandLineBinding{
		ValueFrom: Expression(n.Value),
	}, nil
}

func (l *loader) ScalarToBool(n node) (bool, error) {
	if yamlTrueRX.MatchString(n.Value) {
		return true, nil
	}
	if yamlFalseRX.MatchString(n.Value) {
		return false, nil
	}
	return false, fmt.Errorf("invalid boolean value: %s", n.Value)
}

func (l *loader) ScalarToOptOut(n node) (OptOut, error) {
	if yamlTrueRX.MatchString(n.Value) {
		return OptOut{v: true, set: true}, nil
	}
	if yamlFalseRX.MatchString(n.Value) {
		return OptOut{v: false, set: true}, nil
	}
	return OptOut{}, fmt.Errorf("invalid boolean value: %s", n.Value)
}

func (l *loader) MappingToCommandLineBindingPtr(n node) (*CommandLineBinding, error) {
	clb := CommandLineBinding{}
	err := l.load(n, &clb)
	if err != nil {
		return nil, err
	}
	return &clb, nil
}

func (l *loader) SeqToCommandLineBindingPtrSlice(n node) ([]*CommandLineBinding, error) {
	var clbs []*CommandLineBinding
	for _, c := range n.Children {
		clb := CommandLineBinding{}
		err := l.load(c, &clb)
		if err != nil {
			return nil, err
		}
		clbs = append(clbs, &clb)
	}
	return clbs, nil
}

func (l *loader) SeqToCommandInputSlice(n node) ([]CommandInput, error) {
	var inputs []CommandInput

	for _, c := range n.Children {
		i := CommandInput{}
		err := l.load(c, &i)
		if err != nil {
			return nil, err
		}
		inputs = append(inputs, i)
	}

	return inputs, nil
}

func (l *loader) SeqToCommandOutputSlice(n node) ([]CommandOutput, error) {
	var outputs []CommandOutput

	for _, c := range n.Children {
		i := CommandOutput{}
		err := l.load(c, &i)
		if err != nil {
			return nil, err
		}
		outputs = append(outputs, i)
	}

	return outputs, nil
}

func (l *loader) MappingToCommandInputSlice(n node) ([]CommandInput, error) {
	var inputs []CommandInput

	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v
		i := CommandInput{ID: k}
		if err := l.load(v, &i); err != nil {
			return nil, err
		}
		inputs = append(inputs, i)
	}

	return inputs, nil
}

func (l *loader) MappingToCommandOutputSlice(n node) ([]CommandOutput, error) {
	var outputs []CommandOutput

	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v
		o := CommandOutput{ID: k}
		if err := l.load(v, &o); err != nil {
			return nil, err
		}
		outputs = append(outputs, o)
	}

	return outputs, nil
}

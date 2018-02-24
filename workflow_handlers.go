package cwl

import (
	"fmt"
	"github.com/commondream/yamlast"
)

func (l *loader) MappingToWorkflowInputSlice(n node) ([]WorkflowInput, error) {
	var inputs []WorkflowInput

	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v
		i := WorkflowInput{ID: k}

		switch v.Kind {
		case yamlast.MappingNode:
			err := l.load(v, &i)
			if err != nil {
				return nil, err
			}
		case yamlast.ScalarNode:
			err := l.load(v, &i.Type)
			if err != nil {
				return nil, err
			}
		default:
			// TODO want errors to return position information
			return nil, fmt.Errorf("invalid yaml node type for workflow input")
		}

		inputs = append(inputs, i)
	}

	return inputs, nil
}

func (l *loader) MappingToWorkflowOutputSlice(n node) ([]WorkflowOutput, error) {
	var outputs []WorkflowOutput

	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v
		o := WorkflowOutput{ID: k}

		switch v.Kind {
		case yamlast.MappingNode:
			err := l.load(v, &o)
			if err != nil {
				return nil, err
			}

		case yamlast.ScalarNode:
			err := l.load(v, &o.Type)
			if err != nil {
				return nil, err
			}
		default:
			// TODO want errors to return position information
			return nil, fmt.Errorf("invalid yaml node type for workflow output")
		}
		outputs = append(outputs, o)
	}

	return outputs, nil
}

func (l *loader) MappingToStepSlice(n node) ([]Step, error) {
	steps := []Step{}
	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v
		step := Step{ID: k}
		err := l.load(v, &step)
		if err != nil {
			return nil, err
		}
		steps = append(steps, step)
	}
	return steps, nil
}

func (l *loader) ScalarToStepOutput(n node) (StepOutput, error) {
	return StepOutput{ID: n.Value}, nil
}

func (l *loader) SeqToStepInputSlice(n node) ([]StepInput, error) {
	ins := []StepInput{}
	for _, c := range n.Children {
		in := StepInput{}
		err := l.load(c, &in)
		if err != nil {
			return nil, err
		}

		ins = append(ins, in)
	}
	return ins, nil
}

func (l *loader) SeqToStepOutputSlice(n node) ([]StepOutput, error) {
	outs := []StepOutput{}
	for _, c := range n.Children {
		out := StepOutput{}
		err := l.load(c, &out)
		if err != nil {
			return nil, err
		}
		outs = append(outs, out)
	}
	return outs, nil
}

func (l *loader) MappingToStepInputSlice(n node) ([]StepInput, error) {
	ins := []StepInput{}
	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v
		in := StepInput{ID: k}

		switch v.Kind {
		case yamlast.MappingNode:
			err := l.load(v, &in)
			if err != nil {
				return nil, err
			}

		case yamlast.ScalarNode:
			in.Source = []string{v.Value}
		default:
			// TODO want errors to return position information
			return nil, fmt.Errorf("invalid yaml node type for step input")
		}

		ins = append(ins, in)
	}
	return ins, nil
}

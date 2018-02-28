package cwl

func (l *loader) ScalarToCommandOutput(n node) (CommandOutput, error) {
	o := CommandOutput{}
	err := l.load(n, &o.Type)
	return o, err
}

func (l *loader) ScalarToCommandLineBinding(n node) (CommandLineBinding, error) {
	return CommandLineBinding{
		ValueFrom: Expression(n.Value),
	}, nil
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

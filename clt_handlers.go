package cwl

func loadOutputScalar(l *loader, n node) (interface{}, error) {
	o := CommandOutput{}
	err := l.load(n, &o.Type)
	return o, err
}

func loadBindingScalar(l *loader, n node) (interface{}, error) {
	return CommandLineBinding{
		ValueFrom: Expression(n.Value),
	}, nil
}

func loadInputsMapping(l *loader, n node) (interface{}, error) {
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

func loadOutputsMapping(l *loader, n node) (interface{}, error) {
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

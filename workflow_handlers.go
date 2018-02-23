package cwl

func loadWorkflowInputScalar(l *loader, n node) (interface{}, error) {
  i := WorkflowInput{}
	err := l.load(n, &i.Type)
	return i, err
}

func loadWorkflowOutputScalar(l *loader, n node) (interface{}, error) {
  o := WorkflowOutput{}
	err := l.load(n, &o.Type)
	return o, err
}

func loadWorkflowInputsMapping(l *loader, n node) (interface{}, error) {
	var inputs []WorkflowInput

	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v
		i := WorkflowInput{ID: k}
		if err := l.load(v, &i); err != nil {
			return nil, err
		}
		inputs = append(inputs, i)
	}

	return inputs, nil
}

func loadWorkflowOutputsMapping(l *loader, n node) (interface{}, error) {
	var outputs []WorkflowOutput

	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v
		o := WorkflowOutput{ID: k}
		if err := l.load(v, &o); err != nil {
			return nil, err
		}
		outputs = append(outputs, o)
	}

	return outputs, nil
}

func loadWorkflowStepMapping(l *loader, n node) (interface{}, error) {
  return nil, nil
}

package cwl

import (
	"strings"
)

func (l *loader) SeqToInputValue(n node) (InputValue, error) {
	vals := []InputValue{}
	for _, c := range n.Children {
		var a InputValue
		err := l.load(c, &a)
		if err != nil {
			return nil, err
		}
		vals = append(vals, a)
	}
	return vals, nil
}

func (l *loader) MappingToInputValue(n node) (InputValue, error) {

	class := findKey(n, "class")
	switch strings.ToLower(class) {
	case "file":
		f := File{}
		err := l.load(n, &f)
		if err != nil {
			return nil, err
		}
		return f, nil
	case "directory":
		f := Directory{}
		err := l.load(n, &f)
		if err != nil {
			return nil, err
		}
		return f, nil
	}

	vals := map[string]InputValue{}
	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v

		var a InputValue
		err := l.load(v, &a)
		if err != nil {
			return nil, err
		}
		vals[k] = a
	}
	return vals, nil
}

func (l *loader) MappingToInputValues(n node) (InputValues, error) {
	vals := InputValues{}
	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v
		var a InputValue
		err := l.load(v, &a)
		if err != nil {
			return nil, err
		}
		vals[k] = a
	}
	return vals, nil
}

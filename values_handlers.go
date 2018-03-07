package cwl

import (
	"strings"
)

func (l *loader) SeqToValue(n node) (Value, error) {
	vals := []Value{}
	for _, c := range n.Children {
		var a Value
		err := l.load(c, &a)
		if err != nil {
			return nil, err
		}
		vals = append(vals, a)
	}
	return vals, nil
}

func (l *loader) MappingToValue(n node) (Value, error) {

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

	vals := map[string]Value{}
	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v

		var a Value
		err := l.load(v, &a)
		if err != nil {
			return nil, err
		}
		vals[k] = a
	}
	return vals, nil
}

func (l *loader) MappingToValues(n node) (Values, error) {
	vals := Values{}
	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v
		var a Value
		err := l.load(v, &a)
		if err != nil {
			return nil, err
		}
		vals[k] = a
	}
	return vals, nil
}

package cwl

import (
  "fmt"
  "strings"
)

func loadRequirementsMapping(l *loader, n node) (interface{}, error) {
	var reqs []Requirement
	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v
		x, err := loadHintByName(l, strings.ToLower(k), v)
		if err != nil {
			return nil, err
		}
		req := x.(Requirement)
		reqs = append(reqs, req)
	}
	return reqs, nil
}

func loadHintsMapping(l *loader, n node) (interface{}, error) {
	var hints []Hint
	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v
		h, err := loadHintByName(l, strings.ToLower(k), v)
		if err != nil {
			return nil, err
		}
		hint := h.(Hint)
		hints = append(hints, hint)
	}
	return hints, nil
}

func loadHintMapping(l *loader, n node) (interface{}, error) {
	class := findKey(n, "class")
	return loadHintByName(l, class, n)
}

func loadHintByName(l *loader, name string, n node) (interface{}, error) {
	switch name {
	case "dockerrequirement":
		d := DockerRequirement{}
		err := l.load(n, &d)
		return d, err
	case "resourcerequirement":
		r := ResourceRequirement{}
		err := l.load(n, &r)
		return r, err
	case "inlinejavascriptrequirement":
		j := InlineJavascriptRequirement{}
		err := l.load(n, &j)
		return j, err
	default:
		return nil, fmt.Errorf("unknown hint name: %s", name)
	}
}

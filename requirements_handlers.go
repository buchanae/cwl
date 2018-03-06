package cwl

import (
	"fmt"
	"github.com/commondream/yamlast"
	"strings"
)

func (l *loader) SeqToRequirementSlice(n node) ([]Requirement, error) {
	var reqs []Requirement
	for _, c := range n.Children {
		switch c.Kind {
		case yamlast.MappingNode:
			r, err := l.MappingToRequirement(c)
			if err != nil {
				return nil, err
			}
			reqs = append(reqs, r.(Requirement))
		default:
			panic("unknown node kind")
		}
	}
	return reqs, nil
}

func (l *loader) MappingToRequirementSlice(n node) ([]Requirement, error) {
	var reqs []Requirement
	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v
		x, err := l.loadReqByName(k, v)
		if err != nil {
			return nil, err
		}
		req := x.(Requirement)
		reqs = append(reqs, req)
	}
	return reqs, nil
}

func (l *loader) MappingToHintSlice(n node) ([]Hint, error) {
	var hints []Hint
	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v
		h, err := l.loadReqByName(k, v)
		if err != nil {
			return nil, err
		}
		hint := h.(Hint)
		hints = append(hints, hint)
	}
	return hints, nil
}

func (l *loader) MappingToRequirement(n node) (Requirement, error) {
	class := findKey(n, "class")
	return l.loadReqByName(class, n)
}

func (l *loader) loadReqByName(name string, n node) (Requirement, error) {
	switch strings.ToLower(name) {
	case "dockerrequirement":
		d := DockerRequirement{}
		err := l.load(n, &d)
		return d, err
	case "resourcerequirement":
		r := ResourceRequirement{}
		err := l.load(n, &r)
		return r, err
	case "envvarrequirement":
    // TODO
	case "shellcommandrequirement":
		s := ShellCommandRequirement{}
		err := l.load(n, &s)
		return s, err
	case "inlinejavascriptrequirement":
		j := InlineJavascriptRequirement{}
		err := l.load(n, &j)
		return j, err
	case "schemadefrequirement":
    // TODO
	case "softwarerequirement":
    // TODO
	case "initialworkdirrequirement":
    // TODO
	case "subworkflowfeaturerequirement":
		return SubworkflowFeatureRequirement{}, nil
	case "scatterfeaturerequirement":
		return ScatterFeatureRequirement{}, nil
	case "multipleinputfeaturerequirement":
		return MultipleInputFeatureRequirement{}, nil
	case "stepinputexpressionrequirement":
		return StepInputExpressionRequirement{}, nil
	}
	return nil, fmt.Errorf("unknown hint name: %s", name)
}

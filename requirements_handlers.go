package cwl

import (
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

// TODO this is wrong? test with:
/*
requirements:
  class: ShellCommandRequirement

which should be:

requirements:
  - class: ShellCommandRequirement
*/
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

func (l *loader) MappingToRequirement(n node) (Requirement, error) {
	class := findKey(n, "class")
	return l.loadReqByName(class, n)
}

func (l *loader) SeqToInitialWorkDirListing(n node) (InitialWorkDirListing, error) {
	return InitialWorkDirListing{}, nil
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
		r := EnvVarRequirement{}
		err := l.load(n, &r)
		return r, err
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
		r := SchemaDefRequirement{}
		err := l.load(n, &r)
		return r, err
	case "softwarerequirement":
		r := SoftwareRequirement{}
		err := l.load(n, &r)
		return r, err
	case "initialworkdirrequirement":
		r := InitialWorkDirRequirement{}
		err := l.load(n, &r)
		return r, err
	case "subworkflowfeaturerequirement":
		return SubworkflowFeatureRequirement{}, nil
	case "scatterfeaturerequirement":
		return ScatterFeatureRequirement{}, nil
	case "multipleinputfeaturerequirement":
		return MultipleInputFeatureRequirement{}, nil
	case "stepinputexpressionrequirement":
		return StepInputExpressionRequirement{}, nil
	}
	return UnknownRequirement{Name: name}, nil
	// TODO logging
	//return nil, fmt.Errorf("unknown requirement name: %s", name)
}

package cwl

import (
  "fmt"
  "strings"
	"github.com/commondream/yamlast"
)

func loadRequirementsSeq(l *loader, n node) (interface{}, error) {
	var reqs []Requirement
  for _, c := range n.Children {
    switch c.Kind {
    case yamlast.MappingNode:
      r, err := loadReqMapping(l, c)
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

func loadRequirementsMapping(l *loader, n node) (interface{}, error) {
	var reqs []Requirement
	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v
		x, err := loadReqByName(l, strings.ToLower(k), v)
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
		h, err := loadReqByName(l, strings.ToLower(k), v)
		if err != nil {
			return nil, err
		}
		hint := h.(Hint)
		hints = append(hints, hint)
	}
	return hints, nil
}

func loadReqMapping(l *loader, n node) (interface{}, error) {
	class := findKey(n, "class")
	return loadReqByName(l, class, n)
}

func loadReqByName(l *loader, name string, n node) (interface{}, error) {
	switch name {
	case "dockerrequirement":
		d := DockerRequirement{}
		err := l.load(n, &d)
		return d, err
	case "resourcerequirement":
		r := ResourceRequirement{}
		err := l.load(n, &r)
		return r, err
  case "envvarrequirement":
  case "shellcommandrequirement":
    s := ShellCommandRequirement{}
    err := l.load(n, &s)
    return s, err
	case "inlinejavascriptrequirement":
		j := InlineJavascriptRequirement{}
		err := l.load(n, &j)
		return j, err
  case "schemadefrequirement":
  case "softwarerequirement":
  case "initialworkdirrequirement":
  case "subworkflowfeaturerequirement":
  case "scatterfeaturerequirement":
  case "multipleinputfeaturerequirement":
  case "stepinputexpressionrequirement":
	}
  return nil, fmt.Errorf("unknown hint name: %s", name)
}

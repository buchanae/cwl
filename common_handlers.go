package cwl

import (
	"fmt"
	"github.com/commondream/yamlast"
	"strings"
)

func loadDoc(l *loader, n node) (Document, error) {
	if len(n.Children) != 1 {
		return nil, fmt.Errorf("unexpected document children")
	}

	m := n.Children[0]
	if m.Kind != yamlast.MappingNode {
		return nil, fmt.Errorf("expected mapping node, got: %s", fmtNode(n, ""))
	}

	class := findKey(m, "class")
	switch class {
	case "commandlinetool":
		t := &CommandLineTool{}
		if err := l.load(m, t); err != nil {
			return nil, err
		}
		return t, nil

	case "workflow":
		wf := &Workflow{}
		if err := l.load(m, wf); err != nil {
			return nil, err
		}
		return wf, nil

	default:
		return nil, fmt.Errorf("unknown document class: '%s'", class)
	}
}

func loadExpressionSeq(l *loader, n node) (interface{}, error) {
  return nil, nil
}

func concatStringSeq(l *loader, n node) (interface{}, error) {
  s := ""
  for _, c := range n.Children {
    if c.Kind != yamlast.ScalarNode {
      panic("unhandled seq type")
    }
    s += c.Value
  }
  return s, nil
}

func loadAny(l *loader, n node) (interface{}, error) {
	panic("unhandled Any")
	return nil, nil
}

func loadTypeMappingSlice(l *loader, n node) (interface{}, error) {
	t, err := loadTypeMapping(l, n)
	if err != nil {
		return nil, err
	}
	return []Type{t.(Type)}, nil
}

func loadTypeMapping(l *loader, n node) (interface{}, error) {
	typ := findKey(n, "type")
	switch typ {
	case "array":
		i, ok := findValue(n, "items")
		if !ok {
			panic("")
		}
		a := ArrayType{}
		err := l.load(i, &a.Items)
		return a, err
	case "record":
		return RecordType{}, nil
	case "enum":
		return EnumType{}, nil
	}
	panic("unknown type")
}

func loadTypeScalarSlice(l *loader, n node) (interface{}, error) {

	if strings.HasSuffix(n.Value, "?") {
		name := strings.TrimSuffix(n.Value, "?")
		name = strings.ToLower(name)
		t, ok := TypesByLowercaseName[name]
		if ok {
			return []Type{t.(Type), Null}, nil
		}
	}

	t, err := loadTypeScalar(l, n)
	if err != nil {
		return nil, err
	}
	if t != nil {
		return []Type{t.(Type)}, nil
	}
	panic("unhandled cwl type")
}

// TODO is "string[]?" acceptable?
func loadTypeScalar(l *loader, n node) (interface{}, error) {
	if strings.HasSuffix(n.Value, "[]") {
		name := strings.TrimSuffix(n.Value, "[]")
		if t, ok := TypesByLowercaseName[strings.ToLower(name)]; ok {
			return ArrayType{Items: t}, nil
		}
	}

	name := strings.ToLower(n.Value)
	if t, ok := TypesByLowercaseName[name]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("unhandled scalar type: %s", n.Value)
}

func loadExpressionScalarSlice(l *loader, n node) (interface{}, error) {
	return []Expression{Expression(n.Value)}, nil
}

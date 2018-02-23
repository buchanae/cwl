package cwl

import (
	"fmt"
	"github.com/commondream/yamlast"
	"strings"
)

type DocumentRef struct {
	URL string
}

func (DocumentRef) doctype() {}

func loadDoc(l *loader, n node) (interface{}, error) {
	if n.Kind != yamlast.MappingNode {
		return nil, fmt.Errorf("expected mapping node, got: %s", fmtNode(n, ""))
	}

	class := findKey(n, "class")
	switch class {
	case "commandlinetool":
		t := &CommandLineTool{}
		if err := l.load(n, t); err != nil {
			return nil, err
		}
		return t, nil

	case "workflow":
		wf := &Workflow{}
		if err := l.load(n, wf); err != nil {
			return nil, err
		}
		return wf, nil

	default:
		return nil, fmt.Errorf("unknown document class: '%s'", class)
	}
}

func loadDocumentRef(l *loader, n node) (interface{}, error) {
	return DocumentRef{URL: n.Value}, nil
}

func loadStringSeq(l *loader, n node) (interface{}, error) {
	strs := []string{}
	for _, c := range n.Children {
		strs = append(strs, c.Value)
	}
	return strs, nil
}

func loadCommandLineBindingSeq(l *loader, n node) (interface{}, error) {
	b := []CommandLineBinding{}
	for _, c := range n.Children {
		if c.Kind != yamlast.MappingNode {
			return nil, fmt.Errorf("unhandled command line binding type")
		}
		clb := CommandLineBinding{}
		err := l.load(c, &clb)
		if err != nil {
			return nil, err
		}
		b = append(b, clb)
	}
	return b, nil
}

func concatStringSeq(l *loader, n node) (interface{}, error) {
	s := ""
	for _, c := range n.Children {
		if c.Kind != yamlast.ScalarNode {
			return nil, fmt.Errorf("unhandled string concat type")
		}
		if s != "" {
			s += "\n" + c.Value
		} else {
			s = c.Value
		}
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

func loadExpressionSeq(l *loader, n node) (interface{}, error) {
	exprs := []Expression{}
	for _, c := range n.Children {
		if c.Kind != yamlast.ScalarNode {
			return nil, fmt.Errorf("invalid yaml node type for expression")
		}

		exprs = append(exprs, Expression(c.Value))
	}
	return exprs, nil
}

func loadExpressionScalarSlice(l *loader, n node) (interface{}, error) {
	return []Expression{Expression(n.Value)}, nil
}

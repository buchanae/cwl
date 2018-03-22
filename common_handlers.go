package cwl

import (
	"fmt"
	"github.com/commondream/yamlast"
	"strings"
)

func (l *loader) MappingToDocument(n node) (Document, error) {

  graphNodes, ok := findValue(n, "$graph")
  if ok {
    if graphNodes.Kind != yamlast.SequenceNode {
      return nil, errf("$graph must be a list of objects")
    }
    graph := Graph{}
    err := l.load(n, &graph)
    if err != nil {
      return nil, err
    }
    return graph, nil
  }

	class := findKey(n, "class")
	switch strings.ToLower(class) {

	case "commandlinetool":
		t := &Tool{}
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

	case "expressiontool":
		t := &ExpressionTool{}
		if err := l.load(n, t); err != nil {
			return nil, err
		}
		return t, nil

	default:
		return nil, fmt.Errorf("unknown document class: '%s'", class)
	}
}

func (l *loader) ScalarToDocument(n node) (Document, error) {
	if _, ok := l.resolver.(noResolver); ok {
		return DocumentRef{Location: n.Value}, nil
	}
	b, base, err := l.resolver.Resolve(l.base, n.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve document: %s", err)
	}
	return LoadDocumentBytes(b, base, l.resolver)
}

func (l *loader) ScalarToExpressionSlice(n node) ([]Expression, error) {
	return []Expression{Expression(n.Value)}, nil
}

func (l *loader) MappingToExpressionMap(n node) (map[string]Expression, error) {
	out := map[string]Expression{}
	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v
		expr := Expression("")
		err := l.load(v, &expr)
		if err != nil {
			return nil, errf("loading expression map: %s", err)
		}
		out[k] = expr
	}
	return out, nil
}

func (l *loader) SeqToExpressionMap(n node) (map[string]Expression, error) {
	out := map[string]Expression{}
	for _, c := range n.Children {

		type envdef struct {
			Name  string     `json:"envName"`
			Value Expression `json:"envValue"`
		}

		item := envdef{}
		err := l.load(c, &item)
		if err != nil {
			return nil, errf("loading expression map: %s", err)
		}
		out[item.Name] = item.Value
	}
	return out, nil
}

func (l *loader) SeqToStringSlice(n node) ([]string, error) {
	strs := []string{}
	for _, c := range n.Children {
		strs = append(strs, c.Value)
	}
	return strs, nil
}

func (l *loader) SeqToCommandLineBindingSlice(n node) ([]CommandLineBinding, error) {
	var b []CommandLineBinding
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

func (l *loader) SeqToString(n node) (string, error) {
	s := ""
	for _, c := range n.Children {
		if c.Kind != yamlast.ScalarNode {
			return "", fmt.Errorf("unhandled string concat type")
		}
		if s != "" {
			s += "\n" + c.Value
		} else {
			s = c.Value
		}
	}
	return s, nil
}

func (l *loader) MappingToInputFieldSlice(n node) ([]InputField, error) {
	var fields []InputField

	for _, kv := range itermap(n) {
		k := kv.k
		v := kv.v
		i := InputField{Name: k}
		err := l.load(v, &i)
		if err != nil {
			return nil, err
		}
		fields = append(fields, i)
	}
	return fields, nil
}

/* Type loading is pretty complex.... */

func (l *loader) MappingToInputTypeSlice(n node) ([]InputType, error) {

	typeVal, ok := findValue(n, "type")
	if !ok {
		return nil, fmt.Errorf("missing input type")
	}
	n = transformTypeNode(n)

	var t []InputType
	err := l.load(typeVal, &t)
	if err != nil {
		return nil, err
	}

	if len(t) == 1 {
		switch z := t[0].(type) {
		case InputArray:
			err := l.load(n, &z)
			if err != nil {
				return nil, err
			}
			t[0] = z
		case InputEnum:
			err := l.load(n, &z)
			if err != nil {
				return nil, err
			}
			t[0] = z
		case InputRecord:
			err := l.load(n, &z)
			if err != nil {
				return nil, err
			}
			t[0] = z
		}
	}

	return t, err
}

func (l *loader) MappingToOutputTypeSlice(n node) ([]OutputType, error) {

	typeVal, ok := findValue(n, "type")
	if !ok {
		return nil, fmt.Errorf("missing output type")
	}
	n = transformTypeNode(n)

	var t []OutputType
	err := l.load(typeVal, &t)
	if err != nil {
		return nil, err
	}

	if len(t) == 1 {
		switch z := t[0].(type) {
		case OutputArray:
			err := l.load(n, &z)
			if err != nil {
				return nil, err
			}
			t[0] = z
		case OutputEnum:
			err := l.load(n, &z)
			if err != nil {
				return nil, err
			}
			t[0] = z
		case OutputRecord:
			err := l.load(n, &z)
			if err != nil {
				return nil, err
			}
			t[0] = z
		}
	}

	return t, err
}

func (l *loader) ScalarToInputTypeSlice(n node) ([]InputType, error) {

	n = transformTypeNode(n)

	if n.Kind != yamlast.ScalarNode {
		var out []InputType
		err := l.load(n, &out)
		if err != nil {
			return nil, err
		}
		return out, nil
	}

	t := l.scalarToType(n.Value, true)
	if t == nil {
		return nil, fmt.Errorf("unknown input type: %s", n.Value)
	}

	ot, ok := t.(InputType)
	if !ok {
		return nil, fmt.Errorf("invalid input type: %s", n.Value)
	}
	return []InputType{ot}, nil
}

func (l *loader) ScalarToOutputTypeSlice(n node) ([]OutputType, error) {

	n = transformTypeNode(n)

	if n.Kind != yamlast.ScalarNode {
		var out []OutputType
		err := l.load(n, &out)
		if err != nil {
			return nil, err
		}
		return out, nil
	}

	t := l.scalarToType(n.Value, false)
	if t == nil {
		return nil, fmt.Errorf("unknown output type: %s", n.Value)
	}

	ot, ok := t.(OutputType)
	if !ok {
		return nil, fmt.Errorf("invalid output type: %s", n.Value)
	}
	return []OutputType{ot}, nil
}

func (l *loader) scalarToType(name string, isInput bool) cwltype {

	var t cwltype
	switch strings.ToLower(name) {
	case "":
		return nil
	case "any":
		t = Any{}
	case "null":
		t = Null{}
	case "boolean":
		t = Boolean{}
	case "int":
		t = Int{}
	case "float":
		t = Float{}
	case "long":
		t = Long{}
	case "double":
		t = Double{}
	case "string":
		t = String{}
	case "file":
		t = FileType{}
	case "directory":
		t = DirectoryType{}
	case "stdout":
		t = Stdout{}
	case "stderr":
		t = Stderr{}
	case "record":
		if isInput {
			t = InputRecord{}
		} else {
			t = OutputRecord{}
		}
	case "enum":
		if isInput {
			t = InputEnum{}
		} else {
			t = OutputEnum{}
		}
	case "array":
		if isInput {
			t = InputArray{}
		} else {
			t = OutputArray{}
		}
	default:
		// TODO possibly only create TypeRef for types staring with "#" or otherwise
		//      looking like IRI/URI format?
		return TypeRef{name}
	}

	return t
}

func (l *loader) MappingToSchemaDef(n node) (SchemaDef, error) {
	typeVal, ok := findValue(n, "type")
	if !ok {
		return SchemaDef{}, fmt.Errorf("missing type for schema def")
	}

	name := findKey(n, "name")
	if name == "" {
		return SchemaDef{}, fmt.Errorf("missing name for schema def")
	}

	// TODO support output types?
	switch typeVal.Value {
	case "record":
		t := InputRecord{}
		err := l.load(n, &t)
		return SchemaDef{Name: name, Type: t}, err

	case "array":
		t := InputArray{}
		err := l.load(n, &t)
		return SchemaDef{Name: name, Type: t}, err

	case "enum":
		t := InputEnum{}
		err := l.load(n, &t)
		return SchemaDef{Name: name, Type: t}, err

	default:
		return SchemaDef{}, fmt.Errorf("unknown schema type: %s", typeVal.Value)
	}
}

/* These are here to avoid the automatic loading of slice types in the loader */

func (l *loader) SeqToInputTypeSlice(n node) ([]InputType, error) {
	var out []InputType
	for _, c := range n.Children {
		var t []InputType
		err := l.load(c, &t)
		if err != nil {
			return nil, err
		}
		out = append(out, t...)
	}
	return out, nil
}

func (l *loader) SeqToOutputTypeSlice(n node) ([]OutputType, error) {
	var out []OutputType
	for _, c := range n.Children {
		var t []OutputType
		err := l.load(c, &t)
		if err != nil {
			return nil, err
		}
		out = append(out, t...)
	}
	return out, nil
}

// transformTypeNode handles type name transformations such as "string[]", "string?", etc.
// http://www.commonwl.org/v1.0/Workflow.html#Document_preprocessing
func transformTypeNode(n node) node {
	if n.Kind != yamlast.ScalarNode {
		return n
	}

	//name := strings.ToLower(n.Value)
	name := strings.TrimSpace(n.Value)

	isNullable := false
	isArray := false

	if strings.HasSuffix(name, "?") {
		name = strings.TrimSuffix(name, "?")
		isNullable = true
	}

	if strings.HasSuffix(name, "[]") {
		name = strings.TrimSuffix(name, "[]")
		isArray = true
	}
	n.Value = strings.TrimSpace(name)

	// Copy input node
	out := &yamlast.Node{
		Kind:   n.Kind,
		Line:   n.Line,
		Column: n.Column,
		Value:  n.Value,
	}

	if isArray {
		out = &yamlast.Node{
			Kind:   yamlast.MappingNode,
			Line:   n.Line,
			Column: n.Column,
			Children: []*yamlast.Node{
				{
					Kind:   yamlast.ScalarNode,
					Line:   n.Line,
					Column: n.Column,
					Value:  "type",
				},
				{
					Kind:   yamlast.ScalarNode,
					Line:   n.Line,
					Column: n.Column,
					Value:  "array",
				},
				{
					Kind:   yamlast.ScalarNode,
					Line:   n.Line,
					Column: n.Column,
					Value:  "items",
				},
				out,
			},
		}
	}

	if isNullable {
		out = &yamlast.Node{
			Kind:   yamlast.SequenceNode,
			Line:   n.Line,
			Column: n.Column,
			Children: []*yamlast.Node{
				out,
				{
					Kind:   yamlast.ScalarNode,
					Line:   n.Line,
					Column: n.Column,
					Value:  "null",
				},
			},
		}
	}
	return node(out)
}

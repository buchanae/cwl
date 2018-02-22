package cwl

import (
  "io/ioutil"
  "github.com/commondream/yamlast"
  "github.com/kr/pretty"
  "strings"
  "reflect"
  "fmt"
)

func LoadFile(p string) (Document, error) {
  b, err := ioutil.ReadFile(p)
  if err != nil {
    return nil, err
  }

  // Dump the document string for debugging.
  fmt.Println(string(b))
  return Load(b)
}

func Load(b []byte) (Document, error) {
  // Parse the YAML into an AST
  yamlnode, err := yamlast.Parse(b)
  if err != nil {
    return nil, fmt.Errorf("parsing yaml: %s", err)
  }

  if yamlnode == nil {
    return nil, fmt.Errorf("empty yaml")
  }

  if len(yamlnode.Children) > 1 {
    return nil, fmt.Errorf("unexpected child count")
  }

  // Dump the tree for debugging.
  dump(yamlnode, "")

  // Being recursively processing the tree.
  return loadDoc(l, yamlnode)
}

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

var l = &loader{
  handlers: map[string]handler{
    "mapping -> []cwl.CommandInput": loadInputsMapping,
    "mapping -> []cwl.CommandOutput": loadOutputsMapping,
    "scalar -> cwl.CommandOutput": loadOutputScalar,
    "mapping -> []cwl.Hint": loadHintsMapping,
    "mapping -> cwl.Hint": loadHintMapping,

    "mapping -> []cwl.Requirement": loadRequirementsMapping,
    "mapping -> cwl.Requirement": loadHintMapping,

    "mapping -> cwl.Type": loadTypeMapping,
    "scalar -> cwl.Type": loadTypeScalar,
    "mapping -> []cwl.Type": loadTypeMappingSlice,
    "scalar -> []cwl.Type": loadTypeScalarSlice,

    "mapping -> cwl.Any": loadAny,
    "scalar -> cwl.CommandLineBinding": loadBindingScalar,
    "scalar -> []cwl.Expression": loadExpressionScalarSlice,
  },
}

func loadOutputScalar(l *loader, n node) (interface{}, error) {
  o := CommandOutput{}
  err := l.load(n, &o.Type)
  return o, err
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


func loadBindingScalar(l *loader, n node) (interface{}, error) {
  return CommandLineBinding{
    ValueFrom: Expression(n.Value),
  }, nil
}

func loadInputsMapping(l *loader, n node) (interface{}, error) {
  pretty.Println("LOAD INPUTS", n)
  var inputs []CommandInput

  for _, kv := range itermap(n) {
    k := kv.k
    v := kv.v
    i := CommandInput{ID: k}
    if err := l.load(v, &i); err != nil {
      return nil, err
    }
    inputs = append(inputs, i)
  }

  return inputs, nil
}

func loadOutputsMapping(l *loader, n node) (interface{}, error) {
  pretty.Println("LOAD OUTPUTS", n)
  var outputs []CommandOutput

  for _, kv := range itermap(n) {
    k := kv.k
    v := kv.v
    o := CommandOutput{ID: k}
    if err := l.load(v, &o); err != nil {
      return nil, err
    }
    outputs = append(outputs, o)
  }

  return outputs, nil
}

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


type handler func(l *loader, n node) (interface{}, error)
type loader struct {
  handlers map[string]handler
}

// "t" must be a pointer
func (l *loader) load(n node, t interface{}) error {
  typ := reflect.TypeOf(t).Elem()
  val := reflect.ValueOf(t).Elem()

  nodeKind := "unknown"
  switch n.Kind {
  case yamlast.MappingNode:
    nodeKind = "mapping"
  case yamlast.SequenceNode:
    nodeKind = "sequence"
  case yamlast.ScalarNode:
    nodeKind = "scalar"
  default:
    panic("unknown node kind")
  }
  handlerName := nodeKind + " -> " + typ.String()

  if handler, ok := l.handlers[handlerName]; ok {
    res, err := handler(l, n)
    if err != nil {
      return err
    }
    if res != nil {
      if !reflect.TypeOf(res).AssignableTo(typ) {
        return fmt.Errorf("can't assign value from handler")
      }
      val.Set(reflect.ValueOf(res))
    }
    return nil
  }

  if n.Kind == yamlast.ScalarNode {
    vt := reflect.TypeOf(n.Value)

    if vt.AssignableTo(typ) {
      val.Set(reflect.ValueOf(n.Value))
      return nil
    } else if vt.ConvertibleTo(typ) {
      val.Set(reflect.ValueOf(n.Value).Convert(typ))
      return nil
    } else {
      err := coerceSet(t, n.Value)
      if err == nil {
        return nil
      }
    }
  }

  switch {
  case typ.Kind() == reflect.Struct && n.Kind == yamlast.MappingNode:
    return l.loadMappingToStruct(n, t)
  case typ.Kind() == reflect.Slice && n.Kind == yamlast.SequenceNode:
    for _, c := range n.Children {
      item := reflect.New(typ.Elem())
      err := l.load(c, item.Interface())
      if err != nil {
        return err
      }
      val.Set(reflect.Append(val, item.Elem()))
    }
  default:
    fmt.Println()
    pretty.Println(handlerName, t)
    fmt.Println(fmtNode(n, ""))
    panic("unhandled type")
  }

  return nil
}

// "n" must be a mapping node.
// "t" must be a pointer to a struct.
func (l *loader) loadMappingToStruct(n node, t interface{}) error {

  if n.Kind != yamlast.MappingNode {
    panic("expected mapping node")
  }
  if len(n.Children) % 2 != 0 {
    panic("expected even number of children in mapping")
  }

  typ := reflect.TypeOf(t).Elem()
  val := reflect.ValueOf(t).Elem()
  // track which fields have been set in order to raise an error
  // when a field exists twice.
  already := map[string]bool{}

  for i := 0; i < len(n.Children) - 1; i += 2 {
    k := n.Children[i]
    v := n.Children[i+1]
    name := strings.ToLower(k.Value)

    if _, ok := already[name]; ok {
      panic("duplicate field")
    }
    already[name] = true

    // Find a matching field in the target struct.
    // Names are case insensitive.
    var field reflect.StructField
    var found bool
    for i := 0; i < typ.NumField(); i++ {
      f := typ.Field(i)

      n := f.Name
      if alt, ok := f.Tag.Lookup("cwl"); ok {
        n = alt
      }

      if strings.ToLower(n) == name {
        field = f
        found = true
        break
      }
    }

    if !found {
      continue
    }

    fmt.Printf("%s %s\n", k.Value, v.Value)

    fv := val.FieldByIndex(field.Index)
    err := l.load(v, fv.Addr().Interface())
    if err != nil {
      return err
    }
  }
  return nil
}

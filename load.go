package cwl

import (
  "io/ioutil"
  "github.com/commondream/yamlast"
  "github.com/kr/pretty"
  "strings"
  "reflect"
  "fmt"
)

var l = &loader{
  handlers: map[string]handler{
    "mapping -> []cwl.CommandInput": loadInputsMapping,
    "mapping -> []cwl.CommandOutput": loadOutputsMapping,
    "scalar -> cwl.CommandOutput": loadOutputScalar,
    "mapping -> []cwl.Hint": loadHintsMapping,
    "mapping -> cwl.Hint": loadHintMapping,

    "mapping -> []cwl.Requirement": loadHintsMapping,
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

func loadOutputScalar(l *loader, n node) interface{} {
  o := CommandOutput{}
  l.load(n, &o.Type)
  return o
}

func LoadFile(p string) {
  b, _ := ioutil.ReadFile(p)

  // Dump the document string for debugging.
  fmt.Println(string(b))

  // Parse the YAML into an AST
  yamlnode, err := yamlast.Parse(b)
  if err != nil {
    panic(err)
  }
  // Dump the tree for debugging.
  dump(yamlnode, "")

  if len(yamlnode.Children) > 1 {
    panic("unexpected child count")
  }

  c := yamlnode.Children[0]
  if c.Kind != yamlast.MappingNode {
    panic("expected root document to be a mapping")
  }

  // Being recursively processing the tree.
  d := loadDoc(l, yamlnode)
  pretty.Println(d)
}

func loadAny(l *loader, n node) interface{} {
  return nil
}

func loadTypeMappingSlice(l *loader, n node) interface{} {
  t := loadTypeMapping(l, n)
  return []Type{t.(Type)}
}

func loadTypeMapping(l *loader, n node) interface{} {
  typ := findKey(n, "type")
  switch typ {
  case "array":
    i, ok := findValue(n, "items")
    if !ok {
      panic("")
    }
    a := ArrayType{}
    l.load(i, &a.Items)
    return a
  case "record":
    return RecordType{}
  case "enum":
    return EnumType{}
  }
  panic("")
}

func loadTypeScalarSlice(l *loader, n node) interface{} {
  t := loadTypeScalar(l, n)
  if t != nil {
    return []Type{t.(Type)}
  }
  panic("unhandled cwl type")
}

func loadTypeScalar(l *loader, n node) interface{} {
  if strings.HasSuffix(n.Value, "[]") {
    name := strings.TrimSuffix(n.Value, "[]")
    if t, ok := TypesByLowercaseName[strings.ToLower(name)]; ok {
      return ArrayType{Items: t}
    }
  }

  if t, ok := TypesByLowercaseName[strings.ToLower(n.Value)]; ok {
    return t
  }
  return nil
}

func loadExpressionScalarSlice(l *loader, n node) interface{} {
  return []Expression{Expression(n.Value)}
}

func loadDoc(l *loader, n node) Document {
  if len(n.Children) != 1 {
    panic("")
  }

  m := n.Children[0]
  if m.Kind != yamlast.MappingNode {
    panic(fmt.Errorf("expected mapping node, got: %s", fmtNode(n, "")))
  }

  class := findKey(m, "class")
  switch class {
  case "commandlinetool":
    t := &CommandLineTool{}
    l.load(m, t)
    return t

  default:
    panic(fmt.Errorf("unknown document class: '%s'", class))
  }
}

func loadBindingScalar(l *loader, n node) interface{} {
  return CommandLineBinding{
    ValueFrom: Expression(n.Value),
  }
}

func loadInputsMapping(l *loader, n node) interface{} {
  pretty.Println("LOAD INPUTS", n)
  var inputs []CommandInput

  for k, v := range tomap(n) {
    i := CommandInput{ID: k}
    l.load(v, &i)
    inputs = append(inputs, i)
  }

  return inputs
}

func loadOutputsMapping(l *loader, n node) interface{} {
  pretty.Println("LOAD OUTPUTS", n)
  var outputs []CommandOutput

  for k, v := range tomap(n) {
    o := CommandOutput{ID: k}
    l.load(v, &o)
    outputs = append(outputs, o)
  }

  return outputs
}

func loadHintsMapping(l *loader, n node) interface{} {
  return nil
}

func loadHintMapping(l *loader, n node) interface{} {
  class := findKey(n, "class")
  switch class {
  case "dockerrequirement":
    d := DockerRequirement{}
    l.load(n, &d)
    return d
  case "resourcerequirement":
    r := ResourceRequirement{}
    l.load(n, &r)
    return r
  }
  return nil
}




func tomap(n node) map[string]node {
  if n.Kind != yamlast.MappingNode {
    panic("")
  }
  m := map[string]node{}
  for i := 0; i < len(n.Children) - 1; i += 2 {
    k := n.Children[i]
    v := n.Children[i+1]
    m[k.Value] = v
  }
  return m
}

type handler func(l *loader, n node) interface{}
type loader struct {
  handlers map[string]handler
}

// "t" must be a pointer
func (l *loader) load(n node, t interface{}) {
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
  }
  handlerName := nodeKind + " -> " + typ.String()


  if handler, ok := l.handlers[handlerName]; ok {
    res := handler(l, n)
    if res != nil {
      fmt.Println("LOAD", typ, t, handlerName, res)
      if !reflect.TypeOf(res).AssignableTo(typ) {
        panic(fmt.Errorf("can't assign value from handler"))
      }
      val.Set(reflect.ValueOf(res))
    }
    return
  }

  if n.Kind == yamlast.ScalarNode {
    vt := reflect.TypeOf(n.Value)

    if vt.AssignableTo(typ) {
      val.Set(reflect.ValueOf(n.Value))
      return
    } else if vt.ConvertibleTo(typ) {
      val.Set(reflect.ValueOf(n.Value).Convert(typ))
      return
    } else {
      err := coerceSet(t, n.Value)
      if err == nil {
        return
      }
      fmt.Println("COERCE SET", t, typ, n.Value, err)
    }
  }

  switch {
  case typ.Kind() == reflect.Struct && n.Kind == yamlast.MappingNode:
    l.loadMappingToStruct(n, t)
  case typ.Kind() == reflect.Slice && n.Kind == yamlast.SequenceNode:
    for _, c := range n.Children {
      item := reflect.New(typ.Elem())
      l.load(c, item.Interface())
      val.Set(reflect.Append(val, item.Elem()))
    }
  default:
    fmt.Println()
    pretty.Println(handlerName, t)
    fmt.Println(fmtNode(n, ""))
    panic("")
  }
}

// "n" must be a mapping node.
// "t" must be a pointer to a struct.
func (l *loader) loadMappingToStruct(n node, t interface{}) {
  pretty.Println("LOAD MAPPING", t)

  if n.Kind != yamlast.MappingNode {
    panic("")
  }
  if len(n.Children) % 2 != 0 {
    panic("")
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
    l.load(v, fv.Addr().Interface())
  }
}

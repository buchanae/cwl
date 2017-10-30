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
    "mapping -> []cwl.CommandInput": loadInputs,
    "mapping -> []cwl.CommandOutput": loadOutputs,
    "mapping -> []cwl.Hint": loadHints,
    "mapping -> cwl.InputType": loadInputType,
    "mapping -> cwl.Any": loadAny,
    "mapping -> []cwl.CommandOutputType": loadOutputTypeMapping,
    "scalar -> cwl.CommandLineBinding": loadBindingScalar,
  },
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

func loadInputType(l *loader, n node) interface{} {
  return nil
}

func loadAny(l *loader, n node) interface{} {
  return nil
}

func loadOutputTypeMapping(l *loader, n node) interface{} {
  return nil
}

func loadOutputType(l *loader, n node) interface{} {
  return nil
}

func loadDoc(l *loader, n node) Document {
  if len(n.Children) != 1 {
    panic("")
  }

  m := n.Children[0]
  if m.Kind != yamlast.MappingNode {
    panic(fmt.Errorf("expected mapping node, got: %s", fmtNode(n, "")))
  }

  class := findClass(m)
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

func loadInputs(l *loader, n node) interface{} {
  pretty.Println("LOAD INPUTS", n)
  var inputs []CommandInput

  switch n.Kind {
  case yamlast.MappingNode:
    for k, v := range tomap(n) {
      i := CommandInput{ID: k}
      l.load(v, &i)
      inputs = append(inputs, i)
    }

  case yamlast.SequenceNode:
    for _, c := range n.Children {
      i := CommandInput{}
      l.load(c, &i)
      inputs = append(inputs, i)
    }

  default:
    panic("")
  }

  return inputs
}

func loadOutputs(l *loader, n node) interface{} {
  pretty.Println("LOAD OUTPUTS", n)
  var outputs []CommandOutput

  switch n.Kind {
  case yamlast.MappingNode:
    for k, v := range tomap(n) {
      o := CommandOutput{ID: k}
      l.load(v, &o)
      outputs = append(outputs, o)
    }

  case yamlast.SequenceNode:
    for _, c := range n.Children {
      o := CommandOutput{}
      l.load(c, &o)
      outputs = append(outputs, o)
    }

  default:
    panic("")
  }

  return outputs
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
        panic("")
      }
      val.Set(reflect.ValueOf(res))
    }
    return
  }

  if n.Kind == yamlast.ScalarNode {
    vt := reflect.TypeOf(n.Value)

    if vt.AssignableTo(typ) {
      val.Set(reflect.ValueOf(n.Value))
    } else if vt.ConvertibleTo(typ) {
      val.Set(reflect.ValueOf(n.Value).Convert(typ))
    } else {
      fmt.Println("COERCE SET", typ)
      coerceSet(t, n.Value)
    }
    return
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
    panic("")
  }
}

func loadHints(l *loader, n node) interface{} {
  return nil
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

    fmt.Printf("%s %s\n", k.Value, v.Value)

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

    fv := val.FieldByIndex(field.Index)
    l.load(v, fv.Addr().Interface())
  }
}

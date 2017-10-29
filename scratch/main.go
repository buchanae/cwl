package main

import (
  "io/ioutil"
  "github.com/commondream/yamlast"
  "github.com/kr/pretty"
  "strings"
  "reflect"
  "fmt"
  "os"
)

func main() {
  p := os.Args[1]
  b, _ := ioutil.ReadFile(p)

  // Dump the document string for debugging.
  fmt.Println(string(b))

  // Parse the YAML into an AST.
  yamlnode, err := yamlast.Parse(b)
  if err != nil {
    panic(err)
  }
  n := wrapNode(yamlnode)
  // Dump the tree for debugging.
  dump(yamlnode, "")

  fmt.Println("LOADING")

  d := loadDocument(n)
  pretty.Println(d)
}

func fmtNode(n *yamlast.Node, indent string) string {
  kind := "Unknown"
  switch n.Kind {
  case yamlast.DocumentNode:
    kind = "Document"
  case yamlast.AliasNode:
    kind = "Alias"
  case yamlast.MappingNode:
    kind = "Mapping"
  case yamlast.SequenceNode:
    kind = "Sequence"
  case yamlast.ScalarNode:
    kind = "Scalar"
  }
  return fmt.Sprintf("%-20s Line/col: %3d %3d %40q",
    indent + kind, n.Line + 1, n.Column, n.Value) //, n.Implicit)
}

func dump(n *yamlast.Node, indent string) {
  fmt.Println(fmtNode(n, indent))
  for _, c := range n.Children {
    dump(c, indent + "  ")
  }
}



type Document interface {
  node
  documentType()
}

// isDoc helps mark a type as a Document type succintly.
// It is unexported, so it doesn't add useless fields to the type using it.
type isDoc interface {
  documentType()
}

type Hint interface {
  hintType()
}

type DockerRequirement struct {
  Hint
  node
  DockerPull node
}

type ResourceRequirement struct {
  Hint
  node
  CoresMin node
}


type CommandLineTool struct {
  node
  isDoc

  Class,
  CWLVersion,
  ID item

  BaseCommand,
  Requirements,
  Label,
  Doc,
  Stdin,
  Stdout,
  Stderr,
  SuccessCodes,
  TemporaryFailCodes,
  PermanentFailCodes node

  Hints []Hint
  Arguments []CommandLineBinding
  Inputs []CommandInputParameter
  Outputs []CommandOutputParameter
}

type CommandLineBinding struct {
  node

  Position, Prefix, Separate, ItemSeparator, ValueFrom, ShellQuote, LoadContents node
}

type CommandInputParameter struct {
  node

  ID node // scalar?
  Default node
  Type InputType
  InputBinding CommandLineBinding
}

type CommandOutputParameter struct {
  node

  ID node
  Type OutputType
  OutputBinding CommandOutputBinding
}

type CommandOutputBinding struct {
  node

  Glob, LoadContents, OutputEval node
}



func loadDocument(n node) Document {
  doc := n.(document)
  children := doc.Children()

  if len(children) > 1 {
    panic("unexpected child count")
  }

  m, ok := children[0].(mapping)
  if !ok {
    panic("expected root document to be a mapping")
  }

  switch findClass(m) {
  case "CommandLineTool":
    t := &CommandLineTool{node: n}
    loadContext(m, t)
    return t

  case "Workflow":
  case "ExpressionTool":
  default:
    panic("unknown class")
  }
  return nil
}


func loadContext(n node, c interface{}) {
  switch x := c.(type) {

  case *CommandLineTool:
    m := n.(mapping)
    loadMapping(m, x)

  case *[]Hint:
    switch z := n.(type) {
    case sequence:
      for _, child := range z.Children() {
        loadContext(child, x)
      }

    case mapping:
      switch findClass(z) {
      case "DockerRequirement":
        d := DockerRequirement{node: n}
        loadMapping(z, &d)
        *x = append(*x, d)

      case "ResourceRequirement":
      default:
      }
    }

  case *InputType:
    switch z := n.(type) {
    case scalar:
      name := strings.ToLower(z.Value())
      t, ok := inputTypesByName[name]
      if !ok {
        panic("unknown input type")
      }
      *x = t
    case mapping:
    case sequence:
    }

  case *OutputType:
    switch z := n.(type) {
    case scalar:
      name := strings.ToLower(z.Value())
      t, ok := outputTypesByName[name]
      if !ok {
        panic("unknown output type")
      }
      *x = t
    case mapping:
    case sequence:
    }

  case *[]CommandLineBinding:
    switch z := n.(type) {
    case sequence:
      for _, child := range z.Children() {
        b := CommandLineBinding{node: n}
        loadContext(child, &b)
        *x = append(*x, b)
      }
    }

  case *CommandLineBinding:
    switch z := n.(type) {
    case scalar:
      x.ValueFrom = z
    case mapping:
      loadMapping(z, x)
    case sequence:
    }

  case *[]CommandInputParameter:
    switch z := n.(type) {
    case sequence:
      for _, child := range z.Children() {
        p := CommandInputParameter{node: n}
        loadContext(child, &p)
        *x = append(*x, p)
      }

    case mapping:
    }

  case *CommandInputParameter:
    m := n.(mapping)
    loadMapping(m, x)

  case *[]CommandOutputParameter:
    switch z := n.(type) {
    case sequence:
      for _, child := range z.Children() {
        p := CommandOutputParameter{node: n}
        loadContext(child, &p)
        *x = append(*x, p)
      }

    case mapping:
    }

  case *CommandOutputParameter:
    m := n.(mapping)
    loadMapping(m, x)

  case *CommandOutputBinding:
    m := n.(mapping)
    loadMapping(m, x)

  default:
    panic(fmt.Errorf("unhandled type: %s\n%s", c, n))
  }
}

type InputType interface {
  inputType()
}

type OutputType interface {
  outputType()
}

type NullType struct {
  InputType
  OutputType
}
type IntType struct {
  InputType
  OutputType
}
type FileType struct {
  InputType
  OutputType
}
type DirectoryType struct {
  InputType
  OutputType
}
type ArrayType struct {
  InputType
  Items InputType
}

var (
  Null = NullType{}
  Int = IntType{}
  File = FileType{}
  Directory = DirectoryType{}
)

var inputTypesByName = map[string]InputType{
  "null": Null,
  "int": Int,
  "file": File,
  "directory": Directory,
}
var outputTypesByName = map[string]OutputType{
  "null": Null,
  "int": Int,
  "file": File,
  "directory": Directory,
}


// "dest" must be a pointer to a struct.
func loadMapping(m mapping, dest interface{}) []item {
  var unknown []item

  destType := reflect.TypeOf(dest).Elem()
  destVal := reflect.ValueOf(dest).Elem()
  // track which fields have been set in order to raise an error
  // when a field exists twice.
  already := map[string]bool{}

  for _, item := range m.Items() {
    name := strings.ToLower(item.k.Value())

    if _, ok := already[name]; ok {
      panic("already set field")
    }

    // Find a matching field in the target struct.
    // Names are case insensitive.
    var field reflect.StructField
    var found bool
    for i := 0; i < destType.NumField(); i++ {
      f := destType.Field(i)
      if strings.ToLower(f.Name) == name {
        field = f
        found = true
        break
      }
    }

    if !found {
      unknown = append(unknown, item)
      continue
    }
    already[name] = true

    fieldVal := destVal.FieldByIndex(field.Index)

    srcType := reflect.TypeOf(item)
    if srcType.AssignableTo(field.Type) {
      fieldVal.Set(reflect.ValueOf(item))
      continue
    }

    i := destVal.FieldByIndex(field.Index).Addr().Interface()
    loadContext(item.v, i)
  }
  return unknown
}



type node interface {
  Line() int
  Column() int
}

type item struct {
  k scalar
  v node
}
func (k item) Line() int {
  return k.k.Line()
}
func (k item) Column() int {
  return k.k.Column()
}

type mapping struct {
  wrapper
}
func (m mapping) Items() []item {
  children := m.wrapper.Children

  if len(children) % 2 != 0 {
    panic("expected mapping to have an even number of children")
  }

  var ret []item
  for i := 0; i < len(children); i += 2 {
    k := wrapNode(children[i]).(scalar)
    v := wrapNode(children[i + 1])
    ret = append(ret, item{k, v})
  }
  return ret
}

type sequence struct {
  wrapper
}
func (s sequence) Children() []node {
  var children []node
  for _, c := range s.wrapper.Children {
    children = append(children, wrapNode(c))
  }
  return children
}

type document struct {
  wrapper
}
func (d document) Children() []node {
  var children []node
  for _, c := range d.wrapper.Children {
    children = append(children, wrapNode(c))
  }
  return children
}

type scalar struct {
  wrapper
}
func (s scalar) Value() string {
  return s.wrapper.Value
}

type wrapper struct {
  *yamlast.Node
}
func (n wrapper) Line() int {
  return n.Node.Line + 1
}
func (n wrapper) Column() int {
  return n.Node.Column
}

func findClass(m mapping) string {
  for _, kv := range m.Items() {
    if kv.k.Value() == "class" {
      if val, ok := kv.v.(scalar); ok {
        return val.Value()
      }
    }
  }
  return ""
}

func wrapNode(y *yamlast.Node) node {
  w := wrapper{Node: y}
  switch y.Kind {
  case yamlast.DocumentNode:
    return document{w}
  case yamlast.MappingNode:
    return mapping{w}
  case yamlast.ScalarNode:
    return scalar{w}
  case yamlast.SequenceNode:
    return sequence{w}
  }
  return w
}

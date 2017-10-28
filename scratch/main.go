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
  // Dump the tree for debugging.
  dump(yamlnode, "")

  if len(yamlnode.Children) > 1 {
    panic("unexpected child count")
  }

  c := yamlnode.Children[0]
  if c.Kind != yamlast.MappingNode {
    panic("expected root document to be a mapping")
  }

  // Wrap the node tree in a node type that's easier to work with.
  root := wrapNode(c).(mapping)

  // Being recursively processing the tree.
  doc(root)
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

// Dump the YAML tree for debugging.
func dump(n *yamlast.Node, indent string) {
  fmt.Printf("%s\n", fmtNode(n, indent))
  for _, c := range n.Children {
    dump(c, indent + "  ")
  }
}

// doc determines which type of document this node is:
// - CommandLineTool
// - Workflow
// - ExpressionTool
func doc(m mapping) {
  for _, kv := range m.Items() {
    if kv.k.Value() == "class" {
      val, ok := kv.v.(scalar)
      if !ok {
        panic("unexpected type for document class")
      }

      switch val.Value() {

      case "CommandLineTool":
        t := commandLineTool{}
        loadMapping(m, &t, errorUnknown)
        pretty.Println(t)

      case "Workflow":
      case "ExpressionTool":

      default:
        panic(fmt.Errorf("unexpected document class:", val.Value()))
      }
    }
  }
}






type hint struct {
  Class node
}

func (h *hint) handleMapping(m mapping) {
  loadMapping(m, h, allowUnknown)
}

type hints []hint
func (hs *hints) handleSequence(s sequence) {
  for _, child := range s.Children() {
    h := hint{}
    load(child, &h)
    *hs = append(*hs, h)
  }
}






type commandLineTool struct {
  Class, CWLVersion, BaseCommand,
  ID, Requirements, Label, Doc,
  Stdin, Stdout, Stderr, SuccessCodes, TemporaryFailCodes,
  PermanentFailCodes node

  Arguments commandArguments
  Inputs commandInputs
  Hints hints
  Outputs commandOutputs
}

type commandLineBinding struct {
  Position, Prefix, Separate, ItemSeparator, ValueFrom, ShellQuote, LoadContents node
}
func (c *commandLineBinding) handleScalar(s scalar) {
  c.ValueFrom = s
}
func (c *commandLineBinding) handleMapping(m mapping) {
  loadMapping(m, c, errorUnknown)
}


type commandArguments struct {
  args []commandLineBinding
}
func (c *commandArguments) handleSequence(s sequence) {
  for _, child := range s.Children() {
    arg := commandLineBinding{}
    load(child, &arg)
    c.args = append(c.args, arg)
  }
}

type commandInputParameter struct {
  ID scalar
  Default node
  //Type inputType
  InputBinding commandLineBinding
}

func (c *commandInputParameter) handleMapping(m mapping) {
  loadMapping(m, c, errorUnknown)
}

func (c *commandInputParameter) handleTypeMapping(m mapping) {
}


type inputType interface {
  inputType()
}

type nullType struct {}
func (nullType) inputType() {}

type fileType struct {}
func (fileType) inputType() {}

type arrayType struct {
  Items inputType
}
func (arrayType) inputType() {}



type commandInputs []commandInputParameter
func (c *commandInputs) handleSequence(s sequence) {
  for _, child := range s.Children() {
    i := commandInputParameter{}
    load(child, &i)
    *c = append(*c, i)
  }
}


type commandOutputs struct {}
func (c *commandOutputs) handleSequence(s sequence) {
  for _, child := range s.Children() {
    p := commandOutputParameter{}
    loadMapping(child.(mapping), &p, allowUnknown)
  }
}


type commandOutputParameter struct {
  ID, Type node
  OutputBinding commandOutputBinding
}



type commandOutputBinding struct {
  Glob, LoadContents, OutputEval node
}
func (c *commandOutputBinding) handleMapping(m mapping) {
  loadMapping(m, c, errorUnknown)
}





type mappingHandler interface {
  handleMapping(mapping)
}

type scalarHandler interface {
  handleScalar(scalar)
}

type sequenceHandler interface {
  handleSequence(sequence)
}


type loadMode int
const (
  allowUnknown loadMode = iota
  errorUnknown
)

// "t" must be a pointer to a struct.
func loadMapping(m mapping, t interface{}, mode loadMode) {
  pretty.Println("LOAD", t)

  typ := reflect.TypeOf(t).Elem()
  val := reflect.ValueOf(t).Elem()
  // track which fields have been set in order to raise an error
  // when a field exists twice.
  already := map[string]bool{}

  for _, kv := range m.Items() {
    name := kv.k.Value()

    if _, ok := already[name]; ok {
      panic("already set field")
    }

    // Find a matching field in the target struct.
    // Names are case insensitive.
    tf, found := typ.FieldByNameFunc(func(n string) bool {
      return strings.ToLower(n) == strings.ToLower(name)
    })

    if !found {
      if mode == allowUnknown {
        continue
      }
      panic(fmt.Errorf("unknown field: %s", name))
    }

    i := val.FieldByIndex(tf.Index).Addr().Interface()


    switch x := kv.v.(type) {
    case mapping:
      if h, ok := i.(mappingHandler); ok {
        h.handleMapping(x)
      } else {
        pretty.Println(i, x)
        val.FieldByIndex(tf.Index).Set(reflect.ValueOf(kv))
      }

    case sequence:
      if h, ok := i.(sequenceHandler); ok {
        h.handleSequence(x)
      } else {
        pretty.Println(i, x)
        val.FieldByIndex(tf.Index).Set(reflect.ValueOf(kv))
      }

    case scalar:
      if h, ok := i.(scalarHandler); ok {
        h.handleScalar(x)
      } else {
        val.FieldByIndex(tf.Index).Set(reflect.ValueOf(kv))
      }

    default:
      panic("unhandled node kind")
    }

    already[name] = true
  }
}

func load(n node, i interface{}) {
  switch x := n.(type) {
  case mapping:
    h, ok := i.(mappingHandler)
    if !ok {
      pretty.Println(i, h, ok)
      panic(fmt.Errorf("unhandled mapping type:\n    %s", x))
    }
    h.handleMapping(x)

  case sequence:
    h, ok := i.(sequenceHandler)
    if !ok {
      pretty.Println(i, h, ok)
      panic(fmt.Errorf("unhandled sequence type:\n    %s", x))
    }
    h.handleSequence(x)

  case scalar:
    h, ok := i.(scalarHandler)
    if !ok {
      pretty.Println(i, h, ok)
      panic(fmt.Errorf("unhandled scalar type:\n    %s", x))
    }
    h.handleScalar(x)

  default:
    panic(fmt.Errorf("unhandled type:\n    %s", x))
  }
}




type node interface {
  Line() int
  Column() int
  String() string
}

type keyval struct {
  k scalar
  v node
}
func (k keyval) Line() int {
  return k.k.Line()
}
func (k keyval) Column() int {
  return k.k.Column()
}
func (k keyval) String() string {
  return "TODO"
}

type mapping struct {
  wrapper
}
func (m mapping) Items() []keyval {
  children := m.wrapper.Children

  if len(children) % 2 != 0 {
    panic("expected mapping to have an even number of children")
  }

  var ret []keyval
  for i := 0; i < len(children); i += 2 {
    k := wrapNode(children[i]).(scalar)
    v := wrapNode(children[i + 1])
    ret = append(ret, keyval{k, v})
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
func (n wrapper) String() string {
  return fmtNode(n.Node, "")
}


func wrapNode(y *yamlast.Node) node {
  w := wrapper{y}
  switch y.Kind {
  case yamlast.MappingNode:
    return mapping{w}
  case yamlast.ScalarNode:
    return scalar{w}
  case yamlast.SequenceNode:
    return sequence{w}
  }
  return w
}

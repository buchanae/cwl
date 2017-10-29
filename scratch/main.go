package main

import (
  "io/ioutil"
  "github.com/commondream/yamlast"
  "github.com/ohsu-comp-bio/cwl"
  //"github.com/kr/pretty"
  "strings"
  "fmt"
  "os"
)

func main() {
  p := os.Args[1]
  b, _ := ioutil.ReadFile(p)

  // Dump the document string for debugging.
  fmt.Println(string(b))

  // Parse the YAML into an AST.
  n, err := yamlast.Parse(b)
  if err != nil {
    panic(err)
  }

  x := convert(n, 0)

  loadDocument(x)

  // Dump the tree for debugging.
  dump(x)
}

func convert(n *yamlast.Node, depth int) *node {
  x := &node{
    line: n.Line,
    col: n.Column,
    depth: depth,
  }

  switch n.Kind {
  case yamlast.DocumentNode:
    x.kind = "doc"
    for _, c := range n.Children {
      x.children = append(x.children, convert(c, depth + 1))
    }
    return x

  case yamlast.AliasNode:
    panic("alias node is not handled")

  case yamlast.MappingNode:
    x.kind = "mapping"
    for _, c := range n.Children {
      x.children = append(x.children, convert(c, depth + 1))
    }
    return x

  case yamlast.SequenceNode:
    x.kind = "sequence"
    for _, c := range n.Children {
      x.children = append(x.children, convert(c, depth + 1))
    }
    return x

  case yamlast.ScalarNode:
    x.kind = "scalar"
    x.value = n.Value
    return x

  default:
    panic("unknown yaml node type")
  }
}

func dump(n *node) {
  fmt.Println(n)
  for _, c := range n.children {
    dump(c)
  }
}

type node struct {
  kind string
  line, col, depth int
  value string
  children []*node
}

func (n *node) String() string {
  indent := strings.Repeat("  ", n.depth)
  return fmt.Sprintf("%-20s Line,col: %3d,%3d %40q",
    indent + n.kind, n.line, n.col, n.value)
}

func loadDocument(n *node) {
  if n.kind != "doc" {
    panic("expected document node")
  }
  if len(n.children) != 1 || n.children[0].kind != "mapping" {
    panic("expected document to be mapping")
  }
  mn := n.children[0]

  kindFromClass(mn)

  switch mn.kind {
  case "commandlinetool":
    loadCommandLineTool(mn)
  default:
    panic("unknown class")
  }
}

func loadArgument(n *node) (b cwl.CommandLineBinding) {

  switch n.kind {
  case "mapping":
    for k, v := range tomap(n) {
      switch strings.ToLower(k.value) {
      case "valuefrom":
        b.ValueFrom = cwl.MaybeExpression(v.value)
      case "position":
      case "prefix":
      case "stdout":
      case "stderr":
      }
    }

  case "scalar":
    b.ValueFrom = cwl.MaybeExpression(n.value)

  default:
    panic("unhandled type")
  }
  return
}

func loadArguments(n *node) (args []cwl.CommandLineBinding) {
  switch n.kind {
  case "sequence":
    for _, c := range n.children {
      a := loadArgument(c)
      args = append(args, a)
    }
  default:
    panic("unhandled type")
  }
  return
}

func loadCommandLineTool(n *node) (t cwl.CommandLineTool) {
  for k, v := range tomap(n) {
    switch strings.ToLower(k.value) {
    case "id":
      t.ID = v.value
    case "cwlversion":
      t.CWLVersion = v.value
    case "inputs":
    case "stdin":
      t.Stdin = cwl.MaybeExpression(v.value)
    case "stdout":
      t.Stdout = cwl.MaybeExpression(v.value)
    case "stderr":
      t.Stderr = cwl.MaybeExpression(v.value)
    case "arguments":
      t.Arguments = loadArguments(v)
    case "outputs":
    case "class":
    case "hints":
      switch v.kind {
      case "sequence":
        for _, c := range v.children {
          loadHint(c)
        }
      default:
        panic("unhandled type")
      }
    case "basecommand":
    default:
      panic(fmt.Errorf("unknown field: %s", k.value))
    }
  }
  return
}

func loadHint(n *node) {
  kindFromClass(n)
}

func kindFromClass(n *node) {
  class := findClass(n)
  if class != "" {
    n.kind = class
  }
}

func findClass(n *node) string {
  for i := 0; i < len(n.children) - 1; i += 2 {
    k := n.children[i]
    v := n.children[i+1]
    if strings.ToLower(k.value) == "class" {
      return strings.ToLower(v.value)
    }
  }
  return ""
}


func tomap(n *node) map[*node]*node {
  if len(n.children) % 2 != 0 {
    panic("expected mapping to have even number of children")
  }
  out := map[*node]*node{}
  for i := 0; i < len(n.children); i += 2 {
    k := n.children[i]
    v := n.children[i+1]
    out[k] = v
  }
  return out
}

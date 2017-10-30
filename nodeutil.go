package cwl

import (
  "github.com/commondream/yamlast"
  "strings"
  "fmt"
)

type node *yamlast.Node

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

func findValue(n node, key string) (node, bool) {
  if n.Kind != yamlast.MappingNode {
    panic("")
  }
  for i := 0; i < len(n.Children) - 1; i += 2 {
    k := n.Children[i]
    v := n.Children[i+1]
    if strings.ToLower(k.Value) == strings.ToLower(key) {
      return v, true
    }
  }
  return nil, false
}

func findKey(n node, key string) string {
  if v, ok := findValue(n, key); ok {
    return strings.ToLower(v.Value)
  }
  return ""
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


package cwl

import (
	"github.com/commondream/yamlast"
)

func (l *loader) preprocess(n node) (node, error) {
	switch n.Kind {

	case yamlast.MappingNode:
		for i := 0; i < len(n.Children)-1; i += 2 {
			k := n.Children[i]
			v := n.Children[i+1]
			switch k.Value {
			case "$import":
				b, _, err := l.resolver.Resolve(l.base, v.Value)
				if err != nil {
					return nil, err
				}
				yamlnode, err := yamlast.Parse(b)
				if err != nil {
					return nil, err
				}
				// TODO set line/col/file of the new nodes
				return yamlnode.Children[0], nil

			case "$include":
				b, _, err := l.resolver.Resolve(l.base, v.Value)
				if err != nil {
					return nil, err
				}
				// TODO check line/col of the new node is correct
				return node(&yamlast.Node{
					Kind:   yamlast.ScalarNode,
					Line:   n.Line,
					Column: n.Column,
					Value:  string(b),
				}), nil

			// TODO $mixin

			default:
				x, err := l.preprocess(v)
				if err != nil {
					return nil, err
				}
				n.Children[i+1] = x
			}
		}

	case yamlast.SequenceNode:
		for i, c := range n.Children {
			x, err := l.preprocess(c)
			if err != nil {
				return nil, err
			}
			n.Children[i] = x
		}
	}
	return n, nil
}

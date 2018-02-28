package cwl

import (
	"fmt"
	"github.com/commondream/yamlast"
	"io/ioutil"
)

func LoadFile(p string) (Document, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
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
	//dump(yamlnode, "")

	// Being recursively processing the tree.
	var d Document
	err = l.load(yamlnode.Children[0], &d)
	if err != nil {
		return nil, err
	}
	if d != nil {
		return d, nil
	}
	return nil, nil
}

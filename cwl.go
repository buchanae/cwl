package cwl

import (
	"fmt"
	"github.com/commondream/yamlast"
	"io/ioutil"
)

func Load(loc string) (Document, error) {
	return LoadWithResolver(loc, DefaultResolver{})
}

func LoadWithResolver(loc string, r Resolver) (Document, error) {
	if r == nil {
		r = NoResolve()
	}

	var b []byte
	var base string
	var err error

	// If NoResolve() is being used, load the document bytes using
	// the default resolver, but then continue with NoResolve().
	if _, ok := r.(noResolver); ok {
		d := DefaultResolver{}
		b, base, err = d.Resolve("", loc)
	} else {
		b, base, err = r.Resolve("", loc)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to resolve document: %s", err)
	}
	return LoadDocumentBytes(b, base, r)
}

func LoadDocumentBytes(b []byte, base string, r Resolver) (Document, error) {
	if r == nil {
		r = NoResolve()
	}

	l := loader{base, r}
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

func LoadValuesFile(p string) (Values, error) {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		return nil, err
	}
	return LoadValuesBytes(b)
}

func LoadValuesBytes(b []byte) (Values, error) {
	l := loader{}
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

	v := Values{}
	err = l.load(yamlnode.Children[0], &v)

	if err != nil {
		return nil, err
	}
	return v, nil
}

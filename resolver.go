package cwl

import (
	"io/ioutil"
	"path/filepath"
)

// Resolver describes a type which resolves docment
// by location, such as loading relative file paths
// referenced in a CWL document.
type Resolver interface {
	// Resolve is given the current `base`,
	// such as the directory of the root document,
	// and the `location` referenced by the CWL document.
	//
	// Upon success, the document bytes and the new `base`
	// should be returned.
	Resolve(base, location string) (doc []byte, newBase string, err error)
}

// DefaultResolver is a document location resolver which
// resolves local file paths.
type DefaultResolver struct{}

func (DefaultResolver) Resolve(base, loc string) ([]byte, string, error) {
	if !filepath.IsAbs(loc) {
		loc = filepath.Clean(filepath.Join(base, loc))
	}

	b, err := ioutil.ReadFile(loc)
	if err != nil {
		return nil, "", err
	}
	dir := filepath.Dir(loc)
	return b, dir, nil
}

// NoResolve is a special case resolver which does not
// resolve documents, but instead creates `DocumentRef`
// instances in the document tree.
func NoResolve() Resolver {
	return noResolver{}
}

type noResolver struct {
	Resolver
}

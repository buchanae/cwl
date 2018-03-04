package cwl

import (
	"github.com/buchanae/cwl/fs"
	"testing"
)

func TestResolveFile(t *testing.T) {
	f := File{
		Location: "./examples/record.cwl",
	}

	l := fs.NewLocal()
	x, err := ResolveFile(f, l, true)
	if err != nil {
		t.Error(err)
	}

	debug(x)
}

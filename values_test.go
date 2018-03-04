package cwl

import (
	"github.com/buchanae/cwl/fs"
	"testing"
)

func TestResolveFile(t *testing.T) {
	f := &File{
		Location: "./examples/record.cwl",
	}

	filesys, _ := fs.FindFilesystem(f.Location)
	err := ResolveFile(f, filesys, true)
	if err != nil {
		t.Error(err)
	}

	debug(f)
}

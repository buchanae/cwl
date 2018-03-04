package cwl

import (
	"fmt"
	"github.com/buchanae/cwl/fs"
	"github.com/google/uuid"
	"path/filepath"
	"strings"
)

type InputValue interface{}

type InputValues map[string]InputValue

type File struct {
	Location       string
	Path           string
	Basename       string
	Dirname        string
	Nameroot       string
	Nameext        string
	Checksum       string
	Size           int64
	Format         string
	Contents       string
	SecondaryFiles []Expression
}

type Directory struct {
	Location string
	Path     string
	Basename string
	Listing  []string
}

// ResolveFile uses the filesystem to fill in all fields in the File,
// such as dirname, checksum, size, etc.
func ResolveFile(f File, filesys fs.Filesystem, loadContents bool) (*File, error) {

	// http://www.commonwl.org/v1.0/CommandLineTool.html#File
	// "As a special case, if the path field is provided but the location field is not,
	// an implementation may assign the value of the path field to location,
	// and remove the path field."
	if f.Location == "" && f.Path != "" && f.Contents == "" {
		f.Location = f.Path
		f.Path = ""
	}

	if f.Location == "" && f.Contents == "" {
		return nil, fmt.Errorf("location and contents are empty")
	}

	// If both location and contents are set, one will get overwritten.
	// Can't know which one the caller intended, so fail instead.
	if f.Location != "" && f.Contents != "" {
		return nil, fmt.Errorf("location and contents are both non-empty")
	}

	var x *fs.File
	var err error

	if f.Contents != "" {
		// Determine the file path of the literal.
		// Use the path, or the basename, or generate a random name.
		path := f.Path
		if path == "" {
			path = f.Basename
		}
		if path == "" {
			id, err := uuid.NewRandom()
			if err != nil {
				return nil, fmt.Errorf("failed to generate a random name for a file literal: %s", err)
			}
			path = id.String()
		}

		x, err = filesys.Create(path, f.Contents)
		if err != nil {
			return nil, fmt.Errorf("failed to create file from inline content: %s", err)
		}

	} else {
		x, err = filesys.Info(f.Location)
		if err != nil {
			return nil, fmt.Errorf("failed to get file info: %s", err)
		}

		if loadContents {
			f.Contents, err = filesys.Contents(f.Location)
			if err != nil {
				return nil, fmt.Errorf("failed to load file contents: %s", err)
			}
		}
	}

	f.Location = x.Location
	f.Path = x.Path
	f.Checksum = x.Checksum
	f.Size = x.Size

	// "If basename is provided, it is not required to match the value from location"
	if f.Basename == "" {
		f.Basename = filepath.Base(f.Path)
	}
	f.Nameroot, f.Nameext = splitname(f.Basename)
	f.Dirname = filepath.Dir(f.Path)

	return &f, nil
}

// splitname splits a file name into root and extension,
// with some special CWL rules.
func splitname(n string) (root, ext string) {
	// "For the purposess of path splitting leading periods on the basename are ignored;
	// a basename of .cshrc will have a nameroot of .cshrc."
	x := strings.TrimPrefix(n, ".")
	ext = filepath.Ext(x)
	root = strings.TrimSuffix(n, ext)
	return root, ext
}

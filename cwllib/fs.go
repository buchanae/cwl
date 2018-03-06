package cwllib

import (
	"bytes"
	"fmt"
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
	"strings"
	//"github.com/google/uuid"
	"crypto/sha1"
	"github.com/alecthomas/units"
	"github.com/buchanae/cwl"
)

const maxContentsBytes = 64 * units.Kilobyte

type Filesystem interface {
	Glob(pattern string) ([]*cwl.File, error)

	Create(path, contents string) (*cwl.File, error)
	Info(loc string) (*cwl.File, error)
	Contents(loc string) (string, error)
}

type Local struct {
	workdir string
}

func NewLocal() *Local {
	/*
	  id, err := uuid.NewRandom()
	  if err != nil {
	    return fmt.Errorf("error generating unique file location: %s", err)
	  }
	*/
	return &Local{}
}

func (l *Local) Glob(pattern string) ([]*cwl.File, error) {
	var out []*cwl.File

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	for _, match := range matches {
		f, err := l.Info(match)
		if err != nil {
			return nil, err
		}
		out = append(out, f)
	}
	return out, nil
}

func (l *Local) Create(path, contents string) (*cwl.File, error) {
	if path == "" {
		return nil, fmt.Errorf("can't create file with empty path")
	}

	b := []byte(contents)
	size := int64(len(b))
	if units.MetricBytes(size) > maxContentsBytes {
		return nil, fmt.Errorf("contents is max allowed size (%s)", maxContentsBytes)
	}

	return &cwl.File{
		Location: filepath.Join(l.workdir, path),
		Path:     path,
		Checksum: "sha1$" + fmt.Sprintf("%x", sha1.Sum(b)),
		Size:     size,
	}, nil
}

func (l *Local) Info(loc string) (*cwl.File, error) {
	st, err := os.Stat(loc)
	if err != nil {
		return nil, err
	}

	// TODO make this work with directories
	if st.IsDir() {
		return nil, fmt.Errorf("can't call Info() on a directory: %s", loc)
	}

	return &cwl.File{
		Location: loc,
		Path:     loc,
		// TODO allow config to optionally enable calculating checksum for local files
		Checksum: "",
		Size:     st.Size(),
	}, nil
}

func (l *Local) Contents(loc string) (string, error) {
	fh, err := os.Open(loc)
	if err != nil {
		return "", err
	}
	defer fh.Close()

	buf := &bytes.Buffer{}
	r := &io.LimitedReader{R: fh, N: int64(maxContentsBytes)}
	_, err = io.Copy(buf, r)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ResolveFile uses the filesystem to fill in all fields in the File,
// such as dirname, checksum, size, etc.
func ResolveFile(f cwl.File, filesys Filesystem, loadContents bool) (*cwl.File, error) {

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

	var x *cwl.File
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

	// TODO clean this up. "x" was needed before a package reorg.
	//      possibly can be removed now.
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

func EvalSecondaryFilesPattern(path string, pattern string) string {

	// cwlspec:
	// "If a value in secondaryFiles is a string that is not an expression,
	// it specifies that the following pattern should be applied to the path
	// of the primary file to yield a filename relative to the primary File:"

	// "If string begins with one or more caret ^ characters, for each caret,
	// remove the last file extension from the path (the last period . and all
	// following characters).

	for strings.HasPrefix(pattern, "^") {
		pattern = strings.TrimPrefix(pattern, "^")
		path = strings.TrimSuffix(path, filepath.Ext(path))
	}

	// "Append the remainder of the string to the end of the file path."
	return path + pattern
}

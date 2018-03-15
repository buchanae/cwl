package process

import (
	"errors"
	"github.com/alecthomas/units"
	"github.com/buchanae/cwl"
	"github.com/buchanae/cwl/expr"
	"github.com/google/uuid"
	"path/filepath"
	"strings"
)

var ErrFileNotFound = errors.New("file not found")

type Filesystem interface {
	Create(path, contents string) (*cwl.File, error)
	Info(loc string) (*cwl.File, error)
	Contents(loc string) (string, error)
	Glob(pattern string) ([]*cwl.File, error)
}

const MaxContentsBytes = 64 * units.Kilobyte

// resolveFile uses the filesystem to fill in all fields in the File,
// such as dirname, checksum, size, etc. If f.Contents is given, the
// file will be created via fs.Create(). if `loadContents` is true,
// the file contents will be loaded via fs.Contents().
func (process *Process) resolveFile(f cwl.File, loadContents bool) (*cwl.File, error) {

	// http://www.commonwl.org/v1.0/CommandLineTool.html#File
	// "As a special case, if the path field is provided but the location field is not,
	// an implementation may assign the value of the path field to location,
	// and remove the path field."
	if f.Location == "" && f.Path != "" && f.Contents == "" {
		f.Location = f.Path
		f.Path = ""
	}

	if f.Location == "" && f.Contents == "" {
		return nil, errf("location and contents are empty")
	}

	// If both location and contents are set, one will get overwritten.
	// Can't know which one the caller intended, so fail instead.
	if f.Location != "" && f.Contents != "" {
		return nil, errf("location and contents are both non-empty")
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
				return nil, errf("generating a random name for a file literal: %s", err)
			}
			path = id.String()
		}

		x, err = process.fs.Create(path, f.Contents)
		if err != nil {
			return nil, errf("creating file from inline content: %s", err)
		}

	} else {
		x, err = process.fs.Info(f.Location)
		if err != nil {
			return nil, errf("getting file info for %q: %s", f.Location, err)
		}

		if loadContents {
			f.Contents, err = process.fs.Contents(f.Location)
			if err != nil {
				return nil, errf("loading file contents: %s", err)
			}
		}
	}

	// TODO clean this up. "x" was needed before a package reorg.
	//      possibly can be removed now.
	f.Location = x.Location
	// TODO figure out how to stage files.
	//      namespace inputs so they don't conflict.
	//      remember, the args building depends on this path, so it must happen
	//      in the Process code.
	//f.Path = filepath.Join("/inputs", filepath.Base(x.Path))
	f.Path = filepath.Base(x.Path)
	f.Checksum = x.Checksum
	f.Size = x.Size

	// cwl spec:
	// "If basename is provided, it is not required to match the value from location"
	if f.Basename == "" {
		f.Basename = filepath.Base(f.Path)
	}
	f.Nameroot, f.Nameext = splitname(f.Basename)
	f.Dirname = filepath.Dir(f.Path)

	return &f, nil
}

func (process *Process) resolveSecondaryFiles(file *cwl.File, x cwl.Expression) error {

	// cwl spec:
	// "If the value is an expression, the value of self in the expression
	// must be the primary input or output File object to which this binding applies.
	// The basename, nameroot and nameext fields must be present in self.
	// For CommandLineTool outputs the path field must also be present.
	// The expression must return a filename string relative to the path
	// to the primary File, a File or Directory object with either path
	// or location and basename fields set, or an array consisting of strings
	// or File or Directory objects. It is legal to reference an unchanged File
	// or Directory object taken from input as a secondaryFile.
	// TODO
	if expr.IsExpression(x) {
		process.eval(x, file)
	}

	// cwl spec:
	// "If a value in secondaryFiles is a string that is not an expression,
	// it specifies that the following pattern should be applied to the location
	// of the primary file to yield a filename relative to the primary File:"

	// "If string begins with one or more caret ^ characters, for each caret,
	// remove the last file extension from the location (the last period . and all
	// following characters).
	pattern := string(x)
	// TODO location or path? cwl spec says "path" but I'm suspicious.
	location := file.Location

	for strings.HasPrefix(pattern, "^") {
		pattern = strings.TrimPrefix(pattern, "^")
		location = strings.TrimSuffix(location, filepath.Ext(location))
	}

	// "Append the remainder of the string to the end of the file location."
	sec := cwl.File{
		Location: location + pattern,
	}

	// TODO does LoadContents apply to secondary files? not in the spec
	f, err := process.resolveFile(sec, false)
	if err != nil {
		return err
	}

	file.SecondaryFiles = append(file.SecondaryFiles, f)
	return nil
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

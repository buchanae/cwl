package local

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/alecthomas/units"
	"github.com/buchanae/cwl"
	"github.com/buchanae/cwl/process"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Local struct {
	workdir      string
	CalcChecksum bool
}

func NewLocal(workdir string) *Local {
	return &Local{workdir, false}
}

func (l *Local) Glob(pattern string) ([]cwl.File, error) {
	var out []cwl.File

	pattern = filepath.Join(l.workdir, pattern)

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, errf("%s: %s", err, pattern)
	}

	for _, match := range matches {
		match, _ := filepath.Rel(l.workdir, match)
		f, err := l.Info(match)
		if err != nil {
			return nil, errf("%s: %s", err, match)
		}
		out = append(out, f)
	}
	return out, nil
}

func (l *Local) Create(path, contents string) (cwl.File, error) {
  var x cwl.File
	if path == "" {
		return x, errf("can't create file with empty path")
	}

	b := []byte(contents)
	size := int64(len(b))
	if units.MetricBytes(size) > process.MaxContentsBytes {
		return x, errf("contents is max allowed size (%s)", process.MaxContentsBytes)
	}

	loc := filepath.Join(l.workdir, path)
	abs, err := filepath.Abs(loc)
	if err != nil {
		return x, errf("getting absolute path for %s: %s", loc, err)
	}

	return cwl.File{
		Location: abs,
		Path:     path,
		Checksum: "sha1$" + fmt.Sprintf("%x", sha1.Sum(b)),
		Size:     size,
	}, nil
}

func (l *Local) Info(loc string) (cwl.File, error) {
  var x cwl.File
	if !filepath.IsAbs(loc) {
		loc = filepath.Join(l.workdir, loc)
	}

	st, err := os.Stat(loc)
	if os.IsNotExist(err) {
		return x, process.ErrFileNotFound
	}
	if err != nil {
		return x, err
	}

	// TODO make this work with directories
	if st.IsDir() {
		return x, errf("can't call Info() on a directory: %s", loc)
	}

	abs, err := filepath.Abs(loc)
	if err != nil {
		return x, errf("getting absolute path for %s: %s", loc, err)
	}

	checksum := ""
	if l.CalcChecksum {
		b, err := ioutil.ReadFile(loc)
		if err != nil {
			return x, errf("calculating checksum for %s: %s", loc, err)
		}
		checksum = "sha1$" + fmt.Sprintf("%x", sha1.Sum(b))
	}

	return cwl.File{
		Location: abs,
		Path:     abs,
		Checksum: checksum,
		Size:     st.Size(),
	}, nil
}

func (l *Local) Contents(loc string) (string, error) {
	if !filepath.IsAbs(loc) {
		loc = filepath.Join(l.workdir, loc)
	}

	fh, err := os.Open(loc)
	if os.IsNotExist(err) {
		return "", process.ErrFileNotFound
	}
	if err != nil {
		return "", err
	}
	defer fh.Close()

	buf := &bytes.Buffer{}
	r := &io.LimitedReader{R: fh, N: int64(process.MaxContentsBytes)}
	_, err = io.Copy(buf, r)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

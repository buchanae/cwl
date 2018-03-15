package local

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"crypto/sha1"
	"github.com/alecthomas/units"
	"github.com/buchanae/cwl"
	"github.com/buchanae/cwl/process"
)

type Local struct {
	workdir string
}

func NewLocal(workdir string) *Local {
	return &Local{workdir}
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
		return nil, errf("can't create file with empty path")
	}

	b := []byte(contents)
	size := int64(len(b))
	if units.MetricBytes(size) > process.MaxContentsBytes {
		return nil, errf("contents is max allowed size (%s)", process.MaxContentsBytes)
	}

  loc := filepath.Join(l.workdir, path)
  abs, err := filepath.Abs(loc)
  if err != nil {
    return nil, errf("getting absolute path for %s: %s", loc, err)
  }

	return &cwl.File{
		Location: abs,
		Path:     path,
		Checksum: "sha1$" + fmt.Sprintf("%x", sha1.Sum(b)),
		Size:     size,
	}, nil
}

func (l *Local) Info(loc string) (*cwl.File, error) {
  if !filepath.IsAbs(loc) {
    loc = filepath.Join(l.workdir, loc)
  }

	st, err := os.Stat(loc)
  if os.IsNotExist(err) {
    return nil, process.ErrFileNotFound
  }
	if err != nil {
		return nil, err
	}

	// TODO make this work with directories
	if st.IsDir() {
		return nil, errf("can't call Info() on a directory: %s", loc)
	}

  abs, err := filepath.Abs(loc)
  if err != nil {
    return nil, errf("getting absolute path for %s: %s", loc, err)
  }

	return &cwl.File{
		Location: abs,
		Path:     abs,
		// TODO allow config to optionally enable calculating checksum for local files
		Checksum: "",
		Size:     st.Size(),
	}, nil
}

func (l *Local) Contents(loc string) (string, error) {
  loc = filepath.Join(l.workdir, loc)
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

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
	"github.com/buchanae/cwl/cwllib"
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
	if units.MetricBytes(size) > cwllib.MaxContentsBytes {
		return nil, errf("contents is max allowed size (%s)", cwllib.MaxContentsBytes)
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
		return nil, errf("can't call Info() on a directory: %s", loc)
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
	r := &io.LimitedReader{R: fh, N: int64(cwllib.MaxContentsBytes)}
	_, err = io.Copy(buf, r)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

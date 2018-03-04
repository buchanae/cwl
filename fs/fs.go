package fs

import (
  "bytes"
  "io"
  "fmt"
  "path/filepath"
  "os"
  //"github.com/google/uuid"
  "crypto/sha1"
  "github.com/alecthomas/units"
)

/*
TODO

resolution can be relative, so directory context is required

map location to filesystem provider

*/

type File struct {
  Location string
  Path string
  Checksum string
  Size int64
}

const maxContentsBytes = 64 * units.Kilobyte

type Filesystem interface {

  // TODO secondary files
  // glob
  // check that it exists

  Create(path, contents string) (*File, error)
  Info(loc string) (*File, error)
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

func (l *Local) Create(path, contents string) (*File, error) {
  if path == "" {
    return nil, fmt.Errorf("can't create file with empty path")
  }

  b := []byte(contents)
  size := int64(len(b))
  if units.MetricBytes(size) > maxContentsBytes {
    return nil, fmt.Errorf("contents is max allowed size (%s)", maxContentsBytes)
  }

  return &File{
    Location: filepath.Join(l.workdir, path),
    Path: path,
    Checksum: "sha1$" + fmt.Sprintf("%x", sha1.Sum(b)),
    Size: size,
  }, nil
}

func (l *Local) Info(loc string) (*File, error) {
  st, err := os.Stat(loc)
  if err != nil {
    return nil, err
  }

  // TODO make this work with directories
  if st.IsDir() {
    return nil, fmt.Errorf("can't call Info() on a directory: %s", loc)
  }

  return &File{
    Location: loc,
    Path: loc,
    // TODO allow config to optionally enable calculating checksum for local files
    Checksum: "",
    Size: st.Size(),
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

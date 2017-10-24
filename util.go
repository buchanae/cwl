package cwl

import (
  "fmt"
  "errors"
  "encoding/json"
)

func unmarshal(b []byte) *unmarshaler {
  return &unmarshaler{b:b}
}

type unmarshaler struct {
  b []byte
  done bool
}
func (u *unmarshaler) try(i interface{}) *unmarshaler {
  if !u.done {
    err := json.Unmarshal(u.b, i)
    if err == nil {
      u.done = true
    } else {
      fmt.Println(err)
    }
  }
  return u
}
func (u *unmarshaler) coerce(i interface{}, f func()) *unmarshaler {
  if !u.done {
    err := json.Unmarshal(u.b, i)
    if err == nil {
      f()
      u.done = true
    } else {
      fmt.Println(err)
    }
  }
  return u
}
func (u *unmarshaler) err(s string) error {
  if !u.done {
    return errors.New(s)
  }
  return nil
}

package main

import (
  "github.com/ghodss/yaml"
  "encoding/json"
  "io/ioutil"
  "fmt"
  "reflect"
)

type mapt map[string]interface{}

type Doc struct {
  CWLVersion string
  Class string
}

type Workflow struct {
  Inputs map[string]interface{}
}

type CommandInputParameter struct {
  ID string
  Type string
}


type CommandLineTool struct {
  BaseCommand string
  Inputs CommandLineToolInputs
}

type CommandLineToolInputs map[string]CommandInputParameter
type _CommandLineToolInputs CommandLineToolInputs

func (c CommandLineToolInputs) UnmarshalJSON(b []byte) error {
  _c := make(_CommandLineToolInputs)

  if err := json.Unmarshal(b, &_c); err == nil {
    fmt.Println("map-map", _c)
    return nil
  }

  var arr []CommandInputParameter
  var ms map[string]string

  if err := json.Unmarshal(b, &ms); err == nil {
    for k, v := range ms {
      c[k] = CommandInputParameter{Type: v}
    }
  }

  if err := json.Unmarshal(b, &arr); err == nil {
    for _, i := range arr {
      fmt.Println("param", i)
      c[i.ID] = i
    }
    fmt.Println("err", err, ms == nil)
  }

  return nil
}


/*
func (d *Doc) UnmarshalJSON(b []byte) error {
  json.Unmarshal(b, d)

  fmt.Println(d.CWLVersion)

  return nil
}
*/

func main() {
  b, _ := ioutil.ReadFile("test.yml")

  fmt.Println()
  fmt.Println(string(b))

  doc := Doc{}
  err := yaml.Unmarshal(b, &doc)
  if err != nil {
    panic(err)
  }

  switch doc.Class {
  case "Workflow":
    wf := Workflow{}
    if err := yaml.Unmarshal(b, &wf); err != nil {
      panic(err)
    }
    fmt.Println(wf)
  case "CommandLineTool":
    t := CommandLineTool{}
    if err := yaml.Unmarshal(b, &t); err != nil {
      fmt.Println(reflect.TypeOf(err))
      panic(err)
    }
    fmt.Println(t)
  }
}

package cwl
import (

  "fmt"
  "encoding/json"
)

type ScatterMethod int
const (
  DotProduct ScatterMethod = iota
  NestedCrossProduct
  FlatCrossProduct
)

type LinkMergeMethod int
const (
  Unknown LinkMergeMethod = iota
  MergeNested
  MergeFlattened
)

type CWLType string
const (
  NONE CWLType = "NONE"
  Null = "null"
  Boolean = "boolean"
  Int = "int"
  Long = "long"
  Float = "float"
  Double = "double"
  String = "string"
  FileType = "File"
  DirectoryType = "Directory"
)

func (t *CWLType) UnmarshalJSON(b []byte) error {
  var s string
  err := json.Unmarshal(b, &s)
  if err != nil {
    return err
  }

  switch s {
  case "null":
    *t = Null
  case "boolean":
    *t = Boolean
  case "int":
    *t = Int
  case "long":
    *t = Long
  case "float":
    *t = Float
  case "double":
    *t = Double
  case "string":
    *t = String
  case "File":
    *t = FileType
  case "Directory":
    *t = DirectoryType
  default:
    return fmt.Errorf("Unknown CWLType: %s", s)
  }
  return nil
}

package cwl

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

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

type Type interface {
  cwltype()
}

type isType struct {}
func (isType) cwltype() {}

var (
  Null = isType{}
  Boolean = isType{}
  Int = isType{}
  Float = isType{}
  Long = isType{}
  Double = isType{}
  String = isType{}
  FileType = isType{}
  DirectoryType = isType{}
  Stdout = isType{}
  Stderr = isType{}
)

type RecordType struct {
  isType
}

type NamedType struct {
  Name string
  isType
}

type EnumType struct {
  isType
}

type ArrayType struct {
  Items Type
  isType
}

var TypesByLowercaseName = map[string]Type{
  "null": Null,
  "boolean": Boolean,
  "int": Int,
  "long": Long,
  "float": Float,
  "double": Double,
  "string": String,
  "file": FileType,
  "directory": DirectoryType,
}

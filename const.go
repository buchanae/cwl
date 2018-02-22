package cwl

type Expression string

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

type Primitive string

func (Primitive) cwltype() {}

const (
	Null          = Primitive("null")
	Boolean       = Primitive("boolean")
	Int           = Primitive("int")
	Float         = Primitive("float")
	Long          = Primitive("long")
	Double        = Primitive("double")
	String        = Primitive("string")
	FileType      = Primitive("File")
	DirectoryType = Primitive("Directory")
	Stdout        = Primitive("stdout")
	Stderr        = Primitive("stderr")
)

type RecordType struct {
}

func (RecordType) cwltype() {}

type NamedType struct {
	Name string
}

func (NamedType) cwltype() {}

type EnumType struct {
}

func (EnumType) cwltype() {}

type ArrayType struct {
	Items Type
}

func (ArrayType) cwltype() {}

var TypesByLowercaseName = map[string]Type{
	"null":      Null,
	"boolean":   Boolean,
	"int":       Int,
	"long":      Long,
	"float":     Float,
	"double":    Double,
	"string":    String,
	"file":      FileType,
	"directory": DirectoryType,
	"stdout":    Stdout,
	"stderr":    Stderr,
}

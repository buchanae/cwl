package cwl

type Any interface{}

type Type interface {
	cwltype()
}

type Expression string

type ScatterMethod int

const (
	DotProduct ScatterMethod = iota
	NestedCrossProduct
	FlatCrossProduct
)

// TODO should be a string, so that it serializes nicely
type LinkMergeMethod int

const (
	Unknown LinkMergeMethod = iota
	MergeNested
	MergeFlattened
)

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

type RecordType struct {}

func (RecordType) cwltype() {}

type NamedType struct {
	Name string
}

func (NamedType) cwltype() {}

type EnumType struct {}

func (EnumType) cwltype() {}

type ArrayType struct {
	Items Type
}

func (ArrayType) cwltype() {}

type File struct {
	Location       string
	Path           string
	Basename       string
	Dirname        string
	Nameroot       string
	Nameext        string
	Checksum       string
	Size           int64
	Format         string
	Contents       string
	SecondaryFiles []Expression
}

type Directory struct {
	Location string
	Path     string
	Basename string
	Listing  []string
}

package cwl

import (
	"strings"
)

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

type Null struct{}
type Boolean struct{}
type Int struct{}
type Float struct{}
type Long struct{}
type Double struct{}
type String struct{}
type FileType struct{}
type DirectoryType struct{}
type Stderr struct{}
type Stdout struct{}
type RecordType struct{}
type EnumType struct{}

type ArrayType struct {
	Items Type
}

func (Null) cwltype()                {}
func (Null) String() string          { return "null" }
func (Boolean) cwltype()             {}
func (Boolean) String() string       { return "boolean" }
func (Int) cwltype()                 {}
func (Int) String() string           { return "int" }
func (Float) cwltype()               {}
func (Float) String() string         { return "float" }
func (Long) cwltype()                {}
func (Long) String() string          { return "long" }
func (Double) cwltype()              {}
func (Double) String() string        { return "double" }
func (String) cwltype()              {}
func (String) String() string        { return "string" }
func (FileType) cwltype()            {}
func (FileType) String() string      { return "file" }
func (DirectoryType) cwltype()       {}
func (DirectoryType) String() string { return "directory" }
func (Stderr) cwltype()              {}
func (Stderr) String() string        { return "stderr" }
func (Stdout) cwltype()              {}
func (Stdout) String() string        { return "stdout" }
func (RecordType) cwltype()          {}
func (RecordType) String() string    { return "record" }
func (EnumType) cwltype()            {}
func (EnumType) String() string      { return "enum" }
func (ArrayType) cwltype()           {}
func (ArrayType) String() string     { return "array" }

func GetTypeByName(name string) (Type, bool) {
	switch strings.ToLower(name) {
	case "null":
		return Null{}, true
	case "boolean":
		return Boolean{}, true
	case "int":
		return Int{}, true
	case "long":
		return Long{}, true
	case "float":
		return Float{}, true
	case "double":
		return Double{}, true
	case "string":
		return String{}, true
	case "file":
		return FileType{}, true
	case "directory":
		return DirectoryType{}, true
	case "stdout":
		return Stdout{}, true
	case "stderr":
		return Stderr{}, true
	case "record":
		return RecordType{}, true
	case "array":
		return ArrayType{}, true
	case "enum":
		return EnumType{}, true
	}
	return nil, false
}

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

package cwl

import (
	"fmt"
)

type Expression string

type ScatterMethod string

type Value interface{}
type Values map[string]Value

const (
	DotProduct         ScatterMethod = "dotproduct"
	NestedCrossProduct               = "nested_crossproduct"
	FlatCrossProduct                 = "flat_crossproduct"
)

type LinkMergeMethod string

const (
	MergeNested    LinkMergeMethod = "merge_nested"
	MergeFlattened                 = "merge_flattened"
)

type DocumentRef struct {
	Location string
}

func (d DocumentRef) MarshalText() ([]byte, error) {
	return []byte(d.Location), nil
}

type Any struct{}
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

type FileDir interface {
	filedir()
}

type File struct {
	Location       string    `json:"location,omitempty"`
	Path           string    `json:"path,omitempty"`
	Basename       string    `json:"basename,omitempty"`
	Dirname        string    `json:"dirname,omitempty"`
	Nameroot       string    `json:"nameroot,omitempty"`
	Nameext        string    `json:"nameext,omitempty"`
	Checksum       string    `json:"checksum,omitempty"`
	Size           int64     `json:"size,omitempty"`
	Format         string    `json:"format,omitempty"`
	Contents       string    `json:"contents,omitempty"`
	SecondaryFiles []FileDir `json:"secondaryFiles,omitempty"`
}

type Directory struct {
	Location string    `json:"location,omitempty"`
	Path     string    `json:"path,omitempty"`
	Basename string    `json:"basename,omitempty"`
	Listing  []FileDir `json:"listing,omitempty"`
}

func (File) filedir()      {}
func (Directory) filedir() {}

func (Any) String() string           { return "any" }
func (Null) String() string          { return "null" }
func (Boolean) String() string       { return "boolean" }
func (Int) String() string           { return "int" }
func (Float) String() string         { return "float" }
func (Long) String() string          { return "long" }
func (Double) String() string        { return "double" }
func (String) String() string        { return "string" }
func (FileType) String() string      { return "File" }
func (DirectoryType) String() string { return "Directory" }
func (Stderr) String() string        { return "stderr" }
func (Stdout) String() string        { return "stdout" }
func (InputRecord) String() string   { return "record" }
func (InputEnum) String() string     { return "enum" }
func (InputArray) String() string    { return "array" }
func (OutputRecord) String() string  { return "record" }
func (OutputEnum) String() string    { return "enum" }
func (OutputArray) String() string   { return "array" }

func (Any) MarshalText() ([]byte, error)           { return []byte("any"), nil }
func (Null) MarshalText() ([]byte, error)          { return []byte("null"), nil }
func (Boolean) MarshalText() ([]byte, error)       { return []byte("boolean"), nil }
func (Int) MarshalText() ([]byte, error)           { return []byte("int"), nil }
func (Float) MarshalText() ([]byte, error)         { return []byte("float"), nil }
func (Long) MarshalText() ([]byte, error)          { return []byte("long"), nil }
func (Double) MarshalText() ([]byte, error)        { return []byte("double"), nil }
func (String) MarshalText() ([]byte, error)        { return []byte("string"), nil }
func (FileType) MarshalText() ([]byte, error)      { return []byte("File"), nil }
func (DirectoryType) MarshalText() ([]byte, error) { return []byte("Directory"), nil }
func (Stderr) MarshalText() ([]byte, error)        { return []byte("stderr"), nil }
func (Stdout) MarshalText() ([]byte, error)        { return []byte("stdout"), nil }

type Document interface {
	doctype()
}

func (Tool) doctype()           {}
func (Workflow) doctype()       {}
func (ExpressionTool) doctype() {}
func (DocumentRef) doctype()    {}

type InputType interface {
	String() string
	inputtype()
}

func (Any) inputtype()           {}
func (Null) inputtype()          {}
func (Boolean) inputtype()       {}
func (Int) inputtype()           {}
func (Float) inputtype()         {}
func (Long) inputtype()          {}
func (Double) inputtype()        {}
func (String) inputtype()        {}
func (FileType) inputtype()      {}
func (DirectoryType) inputtype() {}
func (InputRecord) inputtype()   {}
func (InputEnum) inputtype()     {}
func (InputArray) inputtype()    {}

type OutputType interface {
	String() string
	outputtype()
}

func (Any) outputtype()           {}
func (Null) outputtype()          {}
func (Boolean) outputtype()       {}
func (Int) outputtype()           {}
func (Float) outputtype()         {}
func (Long) outputtype()          {}
func (Double) outputtype()        {}
func (String) outputtype()        {}
func (FileType) outputtype()      {}
func (DirectoryType) outputtype() {}
func (Stderr) outputtype()        {}
func (Stdout) outputtype()        {}
func (OutputRecord) outputtype()  {}
func (OutputEnum) outputtype()    {}
func (OutputArray) outputtype()   {}

type cwltype interface {
	String() string
	cwltype()
}

func (Any) cwltype()           {}
func (Null) cwltype()          {}
func (Boolean) cwltype()       {}
func (Int) cwltype()           {}
func (Float) cwltype()         {}
func (Long) cwltype()          {}
func (Double) cwltype()        {}
func (String) cwltype()        {}
func (FileType) cwltype()      {}
func (DirectoryType) cwltype() {}
func (Stderr) cwltype()        {}
func (Stdout) cwltype()        {}
func (InputRecord) cwltype()   {}
func (InputEnum) cwltype()     {}
func (InputArray) cwltype()    {}
func (OutputRecord) cwltype()  {}
func (OutputEnum) cwltype()    {}
func (OutputArray) cwltype()   {}

type Requirement interface {
	requirement()
}

// TODO how many of these could legitimately be used
//      as a hint?
func (UnknownRequirement) requirement()              {}
func (DockerRequirement) requirement()               {}
func (ResourceRequirement) requirement()             {}
func (EnvVarRequirement) requirement()               {}
func (ShellCommandRequirement) requirement()         {}
func (InlineJavascriptRequirement) requirement()     {}
func (SchemaDefRequirement) requirement()            {}
func (SoftwareRequirement) requirement()             {}
func (InitialWorkDirRequirement) requirement()       {}
func (SubworkflowFeatureRequirement) requirement()   {}
func (ScatterFeatureRequirement) requirement()       {}
func (MultipleInputFeatureRequirement) requirement() {}
func (StepInputExpressionRequirement) requirement()  {}

type WorkflowRequirement interface {
	wfrequirement()
}

// OptOut provides a boolean flag that defaults to true.
type OptOut struct {
	v   bool
	set bool
}

func (o *OptOut) Clear() {
	o.v = false
	o.set = false
}

func (o *OptOut) Value() bool {
	if !o.set {
		return true
	}
	return o.v
}

func (o *OptOut) Set(v bool) {
	o.set = true
	o.v = v
}

func (o *OptOut) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%t", o.Value())), nil
}

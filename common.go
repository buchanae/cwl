package cwl

type Any interface{}

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

type DocumentRef struct {
	URL string
}

func (d DocumentRef) MarshalText() ([]byte, error) {
	return []byte(d.URL), nil
}

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

type Document interface {
	doctype()
}

func (CommandLineTool) doctype() {}
func (Workflow) doctype()        {}
func (DocumentRef) doctype()     {}

type InputType interface {
	String() string
	inputtype()
}

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

type Type interface {
	cwltype()
}

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

type Hint interface {
	hint()
}

func (DockerRequirement) hint()               {}
func (ResourceRequirement) hint()             {}
func (EnvVarRequirement) hint()               {}
func (ShellCommandRequirement) hint()         {}
func (InlineJavascriptRequirement) hint()     {}
func (SchemaDefRequirement) hint()            {}
func (SoftwareRequirement) hint()             {}
func (InitialWorkDirRequirement) hint()       {}
func (SubworkflowFeatureRequirement) hint()   {}
func (ScatterFeatureRequirement) hint()       {}
func (MultipleInputFeatureRequirement) hint() {}
func (StepInputExpressionRequirement) hint()  {}

type WorkflowRequirement interface {
	wfrequirement()
}

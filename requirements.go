package cwl

type DockerRequirement struct {
	Pull            string `cwl:"dockerPull"`
	Load            string `cwl:"dockerLoad"`
	File            string `cwl:"dockerFile"`
	Import          string `cwl:"dockerImport"`
	ImageID         string `cwl:"dockerImageID"`
	OutputDirectory string `cwl:"dockerOutputDirectory"`
}

type ResourceRequirement struct {
	CoresMin Expression
	// TODO this is incorrectly denoted in the spec as int | string | expression
	CoresMax  Expression
	RAMMin    Expression
	RAMMax    Expression
	TmpDirMin Expression
	TmpDirMax Expression
	OutDirMin Expression
	OutDirMax Expression
}

type EnvVarRequirement struct {
	Class  string
	EnvDef map[string]Expression
}

type ShellCommandRequirement struct{}

type InlineJavascriptRequirement struct {
	ExpressionLib []string
}

type InputSchema struct{}
type SchemaDefRequirement struct {
	Types []InputSchema
}

type SoftwareRequirement struct {
	Packages []SoftwarePackage
}

type SoftwarePackage struct {
	Package string
	Version []string
	Specs   []string
}

type InitialWorkDirListing struct{}

type InitialWorkDirRequirement struct {
	// TODO the most difficult union type
	Listing InitialWorkDirListing
}

type Dirent struct {
	Entry     Expression
	Entryname Expression
	Writable  bool
}

type SubworkflowFeatureRequirement struct{}

type ScatterFeatureRequirement struct{}

type MultipleInputFeatureRequirement struct{}

type StepInputExpressionRequirement struct{}

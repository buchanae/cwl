package cwl

type Hint interface {
	hint()
}

type Requirement interface {
	requirement()
}
type WorkflowRequirement interface {
	wfrequirement()
}

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

type ShellCommandRequirement struct {}

type InlineJavascriptRequirement struct {
	ExpressionLib []string
}

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

type InitialWorkDirListing struct {}

type InitialWorkDirRequirement struct {
  // TODO the most difficult union type
	Listing InitialWorkDirListing
}

type Dirent struct {
	Entry     Expression
	Entryname Expression
	Writable  bool
}

type SubworkflowFeatureRequirement struct {}

type ScatterFeatureRequirement struct {}

type MultipleInputFeatureRequirement struct {}

type StepInputExpressionRequirement struct {}


// TODO how many of these could legitimately be used
//      as a hint?
func (DockerRequirement) hint()        {}
func (DockerRequirement) requirement() {}
func (ResourceRequirement) hint()        {}
func (ResourceRequirement) requirement() {}
func (EnvVarRequirement) hint()        {}
func (EnvVarRequirement) requirement() {}
func (ShellCommandRequirement) hint()        {}
func (ShellCommandRequirement) requirement() {}
func (InlineJavascriptRequirement) hint()        {}
func (InlineJavascriptRequirement) requirement() {}
func (SchemaDefRequirement) hint()        {}
func (SchemaDefRequirement) requirement() {}
func (SoftwareRequirement) hint()        {}
func (SoftwareRequirement) requirement() {}
func (InitialWorkDirRequirement) hint()        {}
func (InitialWorkDirRequirement) requirement() {}
func (SubworkflowFeatureRequirement) hint()        {}
func (SubworkflowFeatureRequirement) requirement() {}
func (ScatterFeatureRequirement) hint()        {}
func (ScatterFeatureRequirement) requirement() {}
func (MultipleInputFeatureRequirement) hint()        {}
func (MultipleInputFeatureRequirement) requirement() {}
func (StepInputExpressionRequirement) hint()        {}
func (StepInputExpressionRequirement) requirement() {}

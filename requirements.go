package cwl

type DockerRequirement struct {
	Class           string `json:"class"`
	Pull            string `json:"dockerPull"`
	Load            string `json:"dockerLoad"`
	File            string `json:"dockerFile"`
	Import          string `json:"dockerImport"`
	ImageID         string `json:"dockerImageID"`
	OutputDirectory string `json:"dockerOutputDirectory"`
}

type ResourceRequirement struct {
	Class     string     `json:"class"`
	CoresMin  Expression `json:"coresMin"`
	CoresMax  Expression `json:"coresMax"`
	RAMMin    Expression `json:"ramMin"`
	RAMMax    Expression `json:"ramMax"`
	TmpDirMin Expression `json:"tmpdirMin"`
	TmpDirMax Expression `json:"tmpdirMax"`
	OutDirMin Expression `json:"outdirMin"`
	OutDirMax Expression `json:"outdirMax"`
}

type EnvVarRequirement struct {
	Class  string                `json:"class"`
	EnvDef map[string]Expression `json:"envDef"`
}

type ShellCommandRequirement struct {
	Class string `json:"class"`
}

type InlineJavascriptRequirement struct {
	Class         string   `json:"class"`
	ExpressionLib []string `json:"expressionLib"`
}

type InputSchema struct{}
type SchemaDefRequirement struct {
	Class string        `json:"class"`
	Types []InputSchema `json:"types"`
}

type SoftwareRequirement struct {
	Class    string            `json:"class"`
	Packages []SoftwarePackage `json:"packages"`
}

type SoftwarePackage struct {
	Class   string   `json:"class"`
	Package string   `json:"package"`
	Version []string `json:"version"`
	Specs   []string `json:"specs"`
}

type InitialWorkDirListing struct{}

type InitialWorkDirRequirement struct {
	Class string `json:"class"`
	// TODO the most difficult union type
	Listing InitialWorkDirListing `json:"listing"`
}

type Dirent struct {
	Entry     Expression `json:"entry"`
	Entryname Expression `json:"entryname"`
	Writable  bool       `json:"writeable"`
}

type SubworkflowFeatureRequirement struct {
	Class string `json:"class"`
}

type ScatterFeatureRequirement struct {
	Class string `json:"class"`
}

type MultipleInputFeatureRequirement struct {
	Class string `json:"class"`
}

type StepInputExpressionRequirement struct {
	Class string `json:"class"`
}

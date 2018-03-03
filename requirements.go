package cwl

type DockerRequirement struct {
	Class           string `json:"class,omitempty"`
	Pull            string `json:"dockerPull,omitempty"`
	Load            string `json:"dockerLoad,omitempty"`
	File            string `json:"dockerFile,omitempty"`
	Import          string `json:"dockerImport,omitempty"`
	ImageID         string `json:"dockerImageID,omitempty"`
	OutputDirectory string `json:"dockerOutputDirectory,omitempty"`
}

type ResourceRequirement struct {
	Class     string     `json:"class,omitempty"`
	CoresMin  Expression `json:"coresMin,omitempty"`
	CoresMax  Expression `json:"coresMax,omitempty"`
	RAMMin    Expression `json:"ramMin,omitempty"`
	RAMMax    Expression `json:"ramMax,omitempty"`
	TmpDirMin Expression `json:"tmpdirMin,omitempty"`
	TmpDirMax Expression `json:"tmpdirMax,omitempty"`
	OutDirMin Expression `json:"outdirMin,omitempty"`
	OutDirMax Expression `json:"outdirMax,omitempty"`
}

type EnvVarRequirement struct {
	Class  string                `json:"class,omitempty"`
	EnvDef map[string]Expression `json:"envDef,omitempty"`
}

type ShellCommandRequirement struct {
	Class string `json:"class,omitempty"`
}

type InlineJavascriptRequirement struct {
	Class         string   `json:"class,omitempty"`
	ExpressionLib []string `json:"expressionLib,omitempty"`
}

type InputSchema struct{}
type SchemaDefRequirement struct {
	Class string        `json:"class,omitempty"`
	Types []InputSchema `json:"types,omitempty"`
}

type SoftwareRequirement struct {
	Class    string            `json:"class,omitempty"`
	Packages []SoftwarePackage `json:"packages,omitempty"`
}

type SoftwarePackage struct {
	Class   string   `json:"class,omitempty"`
	Package string   `json:"package,omitempty"`
	Version []string `json:"version,omitempty"`
	Specs   []string `json:"specs,omitempty"`
}

type InitialWorkDirListing struct{}

type InitialWorkDirRequirement struct {
	Class string `json:"class,omitempty"`
	// TODO the most difficult union type
	Listing InitialWorkDirListing `json:"listing,omitempty"`
}

type Dirent struct {
	Entry     Expression `json:"entry,omitempty"`
	Entryname Expression `json:"entryname,omitempty"`
	Writable  bool       `json:"writeable,omitempty"`
}

type SubworkflowFeatureRequirement struct {
	Class string `json:"class,omitempty"`
}

type ScatterFeatureRequirement struct {
	Class string `json:"class,omitempty"`
}

type MultipleInputFeatureRequirement struct {
	Class string `json:"class,omitempty"`
}

type StepInputExpressionRequirement struct {
	Class string `json:"class,omitempty"`
}

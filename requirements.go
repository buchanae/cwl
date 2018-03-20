package cwl

type UnknownRequirement struct {
	Name string
}

type DockerRequirement struct {
	Pull            string `json:"dockerPull,omitempty"`
	Load            string `json:"dockerLoad,omitempty"`
	File            string `json:"dockerFile,omitempty"`
	Import          string `json:"dockerImport,omitempty"`
	ImageID         string `json:"dockerImageID,omitempty"`
	OutputDirectory string `json:"dockerOutputDirectory,omitempty"`
}

type ResourceRequirement struct {
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
	EnvDef map[string]Expression `json:"envDef,omitempty"`
}

type ShellCommandRequirement struct {
}

type InlineJavascriptRequirement struct {
	ExpressionLib []string `json:"expressionLib,omitempty"`
}

type SchemaDefRequirement struct {
	Types []SchemaDef `json:"types,omitempty"`
}

type SchemaDef struct {
	Name string `json:"name,omitempty"`
	Type SchemaType
}

type SoftwareRequirement struct {
	Packages []SoftwarePackage `json:"packages,omitempty"`
}

type SoftwarePackage struct {
	Package string   `json:"package,omitempty"`
	Version []string `json:"version,omitempty"`
	Specs   []string `json:"specs,omitempty"`
}

type InitialWorkDirListing struct{}

type InitialWorkDirRequirement struct {
	// TODO the most difficult union type
	Listing InitialWorkDirListing `json:"listing,omitempty"`
}

type Dirent struct {
	Entry     Expression `json:"entry,omitempty"`
	Entryname Expression `json:"entryname,omitempty"`
	Writable  bool       `json:"writeable,omitempty"`
}

type SubworkflowFeatureRequirement struct {
}

type ScatterFeatureRequirement struct {
}

type MultipleInputFeatureRequirement struct {
}

type StepInputExpressionRequirement struct {
}

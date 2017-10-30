package cwl

type Hint struct {}
type Requirement struct {}

type DockerRequirement struct {
  Pull string
  Load string
  File string
  Import string
  ImageID string
  OutputDirectory string
}

type EnvDef struct {}

type EnvVarRequirement struct {
  Class string
  EnvDef EnvDef
}

type EnvironmentDef struct {
  EnvName string
  EnvValue Expression
}

type ShellCommandRequirement struct {
}

type ResourceRequirement struct {
  CoresMin LongExpression
  // TODO this is incorrectly denoted in the spec as int | string | expression
  CoresMax LongExpression
  RAMMin LongExpression
  RAMMax LongExpression
  TmpDirMin LongExpression
  TmpDirMax LongExpression
  OutDirMin LongExpression
  OutDirMax LongExpression
}

type InlineJavascriptRequirement struct {
  ExpressionLib []string
}

type SchemaDefRequirement struct {
  Types []InputSchema
}

type Packages struct {
}

type SoftwareRequirement struct {
  Packages Packages
}

type SoftwarePackage struct {
  Package string
  Version []string
  Specs []string
}

type InitialWorkDirListing struct {
}

type InitialWorkDirRequirement struct {
  Listing InitialWorkDirListing
}

type Dirent struct {
  Entry Expression
  Entryname Expression
  Writable bool
}

type SubworkflowFeatureRequirement struct {
}

type ScatterFeatureRequirement struct {
}

type MultipleInputFeatureRequirement struct {
}

type StepInputExpressionRequirement struct {
}

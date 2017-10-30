package cwl

type Workflow struct {
  Version string `cwl:"cwlVersion"`
  ID string
  Label string
  Doc string

  Hints []Hint
  Requirements []Requirement

  Inputs []WorkflowInput
  Outputs []WorkflowOutput
  Steps []WorkflowStep
}
func (Workflow) doctype() {}

type WorkflowInput interface{}

type Input struct {
  ID string
  Label string
  Doc string
  Streamable bool

  Type []Type
  //InputBinding CommandLineBinding
  Default Any
  SecondaryFiles []Expression
  Format Format
}

type WorkflowOutputs struct {
  ID string
  Type string
  OutputSource string
}

type WorkflowStep struct {
  ID string
  Label string
  Doc string

  Hints []Hint
  Requirements []Requirement

  In []WorkflowStepInput
  Out []WorkflowStepOutput

  //Run WorkflowStepRun
  //Scatter Scatter
  //ScatterMethod ScatterMethod
}

type WorkflowStepOut struct {
}

type WorkflowStepRun struct {
}

type Any interface{}
type Format struct {}
type OutputSources struct {}

type WorkflowOutput struct {
  ID string
  Label string
  Doc string
  Streamable bool

  //OutputBinding CommandOutputBinding
  OutputSource OutputSources
  LinkMerge LinkMergeMethod
  Type []Type
  SecondaryFiles []Expression
  Format Format
}

type Expression string

type File struct {
  Location string
  Path string
  Basename string
  Dirname string
  Nameroot string
  Nameext string
  Checksum string
  Size int64
  Format string
  Contents string
  SecondaryFiles []Expression
}

type Directory struct {
  Location string
  Path string
  Basename string
  Listing []string
}

type OutputRecordSchema struct {
  Type string
  Fields []OutputRecordField
  Label string
}

type OutputRecordFieldType struct {
}

type OutputRecordField struct {
  Name string
  Type OutputRecordFieldType
  Doc string
  //OutputBinding CommandOutputBinding
}

type OutputEnumSchema struct {
  Symbols []string
  Type string
  Label string
  //OutputBinding CommandOutputBinding
}

type OutputArraySchemaItems struct {
  Type
}

type OutputArraySchema struct {
  Items OutputArraySchemaItems
  Type string
  Label string
  //OutputBinding CommandOutputBinding
}



type Scatter struct {
}


type WorkflowStepInputSource struct {
}

type WorkflowStepInput struct {
  ID string
  Source WorkflowStepInputSource
  LinkMerge LinkMergeMethod
  Default Any
  ValueFrom Expression
}

type WorkflowStepOutput struct {
  ID string
}

type InputSchema struct {
}

type InputRecordSchema struct {
  Type string
  Fields []InputRecordField
  Label string
}

type InputRecordFieldType struct {
}

type InputRecordField struct {
  Name string
  Type InputRecordFieldType
  Doc string
  //InputBinding CommandLineBinding
  Label string
}

type InputEnumSchema struct {
  Symbols []string
  Type string
  Label string
  //InputBinding CommandLineBinding
}

type InputArraySchemaItems struct {}

type InputArraySchema struct {
  Items InputArraySchemaItems
  Type string
  Label string
  //InputBinding CommandLineBinding
}

package cwl

type Workflow struct {
  ID string
  Version string
  Label string
  Doc string

  Hints []Hint
  Requirements []Requirement

  Inputs []WorkflowInput
  Outputs []WorkflowOutput
  Steps []WorkflowStep
}

type WorkflowInput interface{}
type WorkflowOutput interface{}

type InputParameter struct {
  ID string
  Label string
  Streamable bool
  Doc string

  Type InputParameterType
  //InputBinding CommandLineBinding
  Default Any
  SecondaryFiles []Expression
  Format ParameterFormat
}

type WorkflowOutputParameters struct {
  ID string
  Type string
  OutputSource string
}

type WorkflowStep struct {
  ID string
  Label string
  Doc string

  In []WorkflowStepInput
  Out []WorkflowStepOutput

  //Run WorkflowStepRun

  Hints []Hint
  Requirements []Requirement
  //Scatter Scatter
  //ScatterMethod ScatterMethod
}

type WorkflowStepOut struct {
}

type WorkflowStepRun struct {
}

type Any interface{}
type ParameterFormat struct {}
type OutputSources struct {}

type OutputParameterType struct {}

type InputParameterArray []InputParameter
// TODO string value here is wrong
type InputParameterMapToType map[string]string
type InputParameterMap map[string]InputParameter

type WorkflowOutputParameterArray []WorkflowOutputParameter
// TODO string value here is wrong
type WorkflowOutputParameterMapToType map[string]string
type WorkflowOutputParameterMap map[string]WorkflowOutputParameter


type WorkflowOutputParameter struct {
  ID string
  Label string
  Streamable bool
  Doc string
  //OutputBinding CommandOutputBinding
  OutputSource OutputSources
  LinkMerge LinkMergeMethod
  Type OutputParameterType
  SecondaryFiles []Expression
  Format ParameterFormat
}

type Expression string

type File struct {
  Class Type
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
  Class Type
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

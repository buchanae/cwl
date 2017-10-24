package cwl

import (
  "fmt"
  "encoding/json"
)

type Workflow struct {
  ID string
  Class string
  Label string
  Doc string
  CWLVersion string

  Inputs map[string]InputParameter
  Outputs map[string]OutputParameter
  Steps []WorkflowStep

  // TODO spec has hints documented incorrectly as array<any>
  //      docker example has it as a map
  //Hints []Hint
  //Requirements []Requirement
}

type InputParameter struct {
  ID string
  Label string
  SecondaryFiles SecondaryFiles
  Format ParameterFormat
  Streamable bool
  Doc Doc
  InputBinding CommandLineBinding
  Default Any
  Type InputParameterType
}

type InputParameters struct {
}

type WorkflowOutputParameters struct {
  ID string
  Type string
  OutputSource string
}

type WorkflowSteps []WorkflowStep

type WorkflowStep struct {
  ID string
  Label string
  Doc string

  In []WorkflowStepInput
  Out []WorkflowStepOutput

  //Run WorkflowStepRun

  Requirements []Requirement
  Hints Hints
  //Scatter Scatter
  //ScatterMethod ScatterMethod
}

type WorkflowStepOut struct {
}

type WorkflowStepRun struct {
  String string
  t *Tool
  wf *Workflow
  et *ExpressionTool
}

func (w *WorkflowStepRun) UnmarshalJSON(b []byte) error {

  switch x := i.(type) {
  case string:
    w.String = x

  case map[string]interface{}:

    switch x["class"] {
    case "workflow":
    case "expression":
    case "commandlinetool"
    default:
      // NOTE can this return useful line/col information? If it can't,
      //      this approach won't be useful.
      return err("unknown workflow step type")
    }
  }
}


type Any interface{}
type Hints interface{}

type Requirement struct {
}

type ParameterFormat struct {
}

type Doc struct {
}

type OutputSources struct {
}

type OutputParameterType struct {
  cwlType CWLType
  stdout bool
  stderr bool
  str string
  arrType []OutputParameterType
  recSchema CommandOutputRecordSchema
  enumSchema CommandOutputEnumSchema
  arrSchema CommandOutputArraySchema
}
func (o *OutputParameterType) UnmarshalJSON(b []byte) error {
  var err error

  err = json.Unmarshal(b, &o.cwlType)
  if err == nil {
    return nil
  }

  var s string
  err = json.Unmarshal(b, &s)
  if err == nil {
    // TODO endswith '?'
    // TODO endswith '[]'
    o.str = s
    return nil
  }

  err = json.Unmarshal(b, &o.arrType)
  if err == nil {
    return nil
  }
  err = json.Unmarshal(b, &o.recSchema)
  if err == nil {
    return nil
  }
  err = json.Unmarshal(b, &o.enumSchema)
  if err == nil {
    return nil
  }
  err = json.Unmarshal(b, &o.arrSchema)
  if err == nil {
    return nil
  }

  return fmt.Errorf("Can't unmarshal OutputParameterType: %s", err)
}


type InputParameterArray []InputParameter
// TODO string value here is wrong
type InputParameterMapToType map[string]string
type InputParameterMap map[string]InputParameter

type WorkflowOutputParameterArray []WorkflowOutputParameter
// TODO string value here is wrong
type WorkflowOutputParameterMapToType map[string]string
type WorkflowOutputParameterMap map[string]WorkflowOutputParameter
/**********************************************************************
 * End shims.
 */


type WorkflowOutputParameter struct {
  ID string
  Label string
  SecondaryFiles SecondaryFiles
  Format ParameterFormat
  Streamable bool
  Doc Doc
  OutputBinding CommandOutputBinding
  OutputSource OutputSources
  LinkMerge LinkMergeMethod
  Type OutputParameterType
}

type Expression string
type MaybeExpression string

type LongMaybeExpression struct {
}

// TODO should the spec be array<string | Expression> ?
// string | Expression | array<string>
type Glob struct {
  arr []string
}
func (g *Glob) UnmarshalJSON(b []byte) error {
  var err error

  var s string
  err = json.Unmarshal(b, &s)
  if err == nil {
    g.arr = append(g.arr, s)
    return nil
  }

  err = json.Unmarshal(b, &g.arr)
  if err == nil {
    return nil
  }

  return err
}

type CommandOutputBinding struct {
  Glob Glob
  LoadContents bool
  OutputEval MaybeExpression
}

type File struct {
  Class CWLType
  Location string
  Path string
  Basename string
  Dirname string
  Nameroot string
  Nameext string
  Checksum string
  Size int64
  SecondaryFiles SecondaryFiles
  Format string
  Contents string
}

type Directory struct {
  Class CWLType
  Location string
  Path string
  Basename string
  Listing SecondaryFiles
}

type OutputRecordSchema struct {
  Type string
  Fields []OutputRecordField
  Label string
}

// CWLType | OutputRecordSchema | OutputEnumSchema | OutputArraySchema | string | array<CWLType | OutputRecordSchema | OutputEnumSchema | OutputArraySchema | string> 
type OutputRecordFieldType struct {
}

type OutputRecordField struct {
  Name string
  Type OutputRecordFieldType
  // TODO why isn't this Doc?
  Doc string
  OutputBinding CommandOutputBinding
}

type OutputEnumSchema struct {
  Symbols []string
  Type string
  Label string
  OutputBinding CommandOutputBinding
}

// CWLType | OutputRecordSchema | OutputEnumSchema | OutputArraySchema | string | array<CWLType | OutputRecordSchema | OutputEnumSchema | OutputArraySchema | string>
type OutputArraySchemaItems struct {
  cwlType CWLType
}
func (o *OutputArraySchemaItems) UnmarshalJSON(b []byte) error {
  var err error
  err = json.Unmarshal(b, &o.cwlType)
  if err == nil {
    return nil
  }
  return fmt.Errorf("OutputArraySchemaItems: %s", err)
}

type OutputArraySchema struct {
  Items OutputArraySchemaItems
  Type string
  Label string
  OutputBinding CommandOutputBinding
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
  ValueFrom MaybeExpression
}

type WorkflowStepOutput struct {
  ID string
}

type InlineJavascriptRequirement struct {
  Class string
  ExpressionLib []string
}

type SchemaDefRequirement struct {
  Class string
  Types []InputSchema
}

type InputSchema struct {
}

type InputRecordSchema struct {
  Type string
  Fields []InputRecordField
  Label string
}

// CWLType | InputRecordSchema | InputEnumSchema | InputArraySchema | string | array<CWLType | InputRecordSchema | InputEnumSchema | InputArraySchema | string>
type InputRecordFieldType struct {
}

type InputRecordField struct {
  Name string
  Type InputRecordFieldType
  Doc string
  InputBinding CommandLineBinding
  Label string
}

type InputEnumSchema struct {
  Symbols []string
  Type string
  Label string
  InputBinding CommandLineBinding
}

type CommandLineBinding struct {
  LoadContents bool
  Position int
  Prefix string
  Separate bool
  ItemSeparator string
  ValueFrom MaybeExpression
  ShellQuote bool
}

// CWLType | InputRecordSchema | InputEnumSchema | InputArraySchema | string | array<CWLType | InputRecordSchema | InputEnumSchema | InputArraySchema | string>
type InputArraySchemaItems struct {
  str string
  cwlType CWLType
}
func (i *InputArraySchemaItems) UnmarshalJSON(b []byte) error {
  var err error
  err = json.Unmarshal(b, &i.cwlType)
  if err == nil {
    return nil
  }

  err = json.Unmarshal(b, &i.str)
  if err == nil {
    return nil
  }
  return err
}

type InputArraySchema struct {
  Items InputArraySchemaItems
  Type string
  Label string
  InputBinding CommandLineBinding
}

type Packages struct {
}

type SoftwareRequirement struct {
  Class string
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
  Class string
  Listing InitialWorkDirListing
}

type Dirent struct {
  Entry MaybeExpression
  Entryname MaybeExpression
  Writable bool
}

type SubworkflowFeatureRequirement struct {
  Class string
}

type ScatterFeatureRequirement struct {
  Class string
}

type MultipleInputFeatureRequirement struct {
  Class string
}

type StepInputExpressionRequirement struct {
  Class string
}

// array<InputParameter> | map<InputParameter.id, InputParameter.type> | map<InputParameter.id, InputParameter>
type ExpressionToolInputs struct {
}

// array<ExpressionToolOutputParameter> | map<ExpressionToolOutputParameter.id, ExpressionToolOutputParameter.type> | map<ExpressionToolOutputParameter.id, ExpressionToolOutputParameter>
type ExpressionToolOutputs struct {
}

type ExpressionTool struct {
  Inputs ExpressionToolInputs
  Outputs ExpressionToolOutputs
  Class string
  Expression MaybeExpression
  ID string
  Requirements []Requirement
  Hints Hints
  Label string
  Doc string
  CWLVersion string
}

// CWLType | CommandInputRecordSchema | CommandInputEnumSchema | CommandInputArraySchema | string | array<CWLType | CommandInputRecordSchema | CommandInputEnumSchema | CommandInputArraySchema | string>
type InputParameterType struct {
  str string
  cwlType CWLType
  arrType CommandInputArraySchema
}
func (i *InputParameterType) UnmarshalJSON(b []byte) error {
  var err error

  err = json.Unmarshal(b, &i.cwlType)
  if err == nil {
    return nil
  }

  var s string
  err = json.Unmarshal(b, &s)
  if err == nil {
    // TODO endswith '?'
    // TODO endswith '[]'
    i.str = s
    return nil
  }

  err = json.Unmarshal(b, &i.arrType)
  if err == nil {
    return nil
  }

  return err
}


type ExpressionToolOutputParameter struct {
  ID string
  Label string
  SecondaryFiles SecondaryFiles
  Format ParameterFormat
  Streamable bool
  Doc Doc
  OutputBinding CommandOutputBinding
  Type OutputParameterType
}

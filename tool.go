package cwl

import (
  "fmt"
  "encoding/json"
)


type Tool struct {
  Inputs ToolInputs
  Outputs ToolOutputs
  Class string
  ID string
  Requirements []Requirement
  Hints Hints
  Label string
  Doc string
  CWLVersion string
  BaseCommand BaseCommand `yaml:"baseCommand"`
  Arguments []Argument
  Stdin MaybeExpression
  Stderr MaybeExpression
  Stdout MaybeExpression
  SuccessCodes []int
  TemporaryFailCodes []int
  PermanentFailCodes []int
}


// array<CommandInputParameter> | map<CommandInputParameter.id, CommandInputParameter.type> | map<CommandInputParameter.id, CommandInputParameter>
type ToolInputs []InputParameter 
func (c *ToolInputs) UnmarshalJSON(b []byte) error {
  var err error

  // Try unmarshaling a map-type
  m := map[string]InputParameter{}
  err = json.Unmarshal(b, &m)
  if err == nil {
    for k, v := range m {
      v.ID = k
      *c = append(*c, v)
    }
    return nil
  }

  // Try unmarshaling a list-type
  a := []InputParameter{}
  err = json.Unmarshal(b, &a)
  if err == nil {
    *c = append(*c, a...)
    return nil
  }

  return fmt.Errorf("Can't unmarshal ToolInputs: %s", err)
}

// array<CommandOutputParameter> | map<CommandOutputParameter.id, CommandOutputParameter.type> | map<CommandOutputParameter.id, CommandOutputParameter>
type ToolOutputs []CommandOutputParameter

func (c *ToolOutputs) UnmarshalJSON(b []byte) error {
  var err error

  m := map[string]CommandOutputParameter{}
  err = json.Unmarshal(b, &m)
  if err == nil {
    for k, v := range m {
      v.ID = k
      *c = append(*c, v)
    }
    return nil
  }

  a := []CommandOutputParameter{}
  err = json.Unmarshal(b, &a)
  if err == nil {
    *c = append(*c, a...)
    return nil
  }

  return fmt.Errorf("Can't unmarshal ToolOutputs: %s", err)
}

// string | array<string>
type BaseCommand []string
func (bcmd *BaseCommand) UnmarshalJSON(b []byte) error {
  var err error

  var s string
  err = json.Unmarshal(b, &s)
  if err == nil {
    *bcmd = append(*bcmd, s)
    return nil
  }
/*
  err = json.Unmarshal(b, &bcmd)
  if err == nil {
    return nil
  }
  */

  return err
}

// string | Expression | CommandLineBinding
type Argument struct {
  str string
  clb CommandLineBinding
}
func (a *Argument) UnmarshalJSON(b []byte) error {
  var err error
  err = json.Unmarshal(b, &a.str)

  if err == nil {
    return nil
  }

  err = json.Unmarshal(b, &a.clb)
  if err == nil {
    return nil
  }

  return err
}
func (a *Argument) String() string {
  // TODO I think evaluation would happen here
  return a.str
}



type CommandInputParameter InputParameter
// CWLType | CommandInputRecordSchema | CommandInputEnumSchema | CommandInputArraySchema | string | array<CWLType | CommandInputRecordSchema | CommandInputEnumSchema | CommandInputArraySchema | string>
type CommandInputParameterType InputParameterType
type CommandInputRecordSchema InputRecordSchema
type CommandInputRecordField InputRecordField
type CommandInputEnumSchema InputEnumSchema
type CommandInputArraySchema InputArraySchema

type CommandOutputRecordSchema OutputRecordSchema
type CommandOutputRecordField OutputRecordField
type CommandOutputEnumSchema OutputEnumSchema
type CommandOutputArraySchema OutputArraySchema

type CommandOutputParameter struct {
  ID string
  Label string
  SecondaryFiles SecondaryFiles
  Format ParameterFormat
  Streamable bool
  Doc Doc
  OutputBinding CommandOutputBinding
  // TODO this can be an array?
  Type OutputParameterType
}

type DockerRequirement struct {
  Class string
  DockerPull string
  DockerLoad string
  DockerFile string
  DockerImport string
  DockerImageID string
  DockerOutputDirectory string
}

// array<EnvironmentDef> | map<EnvironmentDef.envName, EnvironmentDef.envValue> | map<EnvironmentDef.envName, EnvironmentDef>
type EnvDef struct {
}

type EnvVarRequirement struct {
  Class string
  EnvDef EnvDef
}

type EnvironmentDef struct {
  EnvName string
  EnvValue MaybeExpression
}

type ShellCommandRequirement struct {
  Class string
}

type ResourceRequirement struct {
  Class string
  CoresMin LongMaybeExpression
  // TODO this is incorrectly denoted in the spec as int | string | expression
  CoresMax LongMaybeExpression
  RAMMin LongMaybeExpression
  RAMMax LongMaybeExpression
  TmpDirMin LongMaybeExpression
  TmpDirMax LongMaybeExpression
  OutDirMin LongMaybeExpression
  OutDirMax LongMaybeExpression
}

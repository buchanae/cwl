package cwl

type Document interface {
  doctype()
}
type isDoc struct {}
func (isDoc) doctype() {}

type CommandLineTool struct {
  Version string `cwl:"cwlVersion"`
  ID string
  Label string
  Doc string

  Hints []Hint
  Requirements []Requirement

  Inputs []CommandInput
  Outputs []CommandOutput

  BaseCommand []string
  Arguments []CommandLineBinding

  Stdin  Expression
  Stderr Expression
  Stdout Expression

  SuccessCodes []int
  TemporaryFailCodes []int
  PermanentFailCodes []int

  isDoc
}

type CommandLineBinding struct {
  LoadContents bool
  Position int
  Prefix string
  Separate bool
  ItemSeparator string
  ValueFrom Expression
  ShellQuote bool
}

type InputType interface {}

type CommandInput struct {
  ID string
  Label string
  Doc string

  InputBinding CommandLineBinding
  Default Any
  Type InputType
}

type CommandOutputBinding struct {
  Glob []Expression
  LoadContents bool
  OutputEval Expression
}

type CommandOutputType interface {}

type CommandOutput struct {
  ID string
  Label string
  Doc string
  Streamable bool

  SecondaryFiles []Expression
  Formats []Expression

  OutputBinding CommandOutputBinding
  Type []CommandOutputType
}

package cwl

type ExpressionTool struct {
  ID string
  Label string
  Doc string
  Version string

  Hints []Hint
  Requirements []Requirement

  Inputs ExpressionToolInputs
  Outputs ExpressionToolOutputs

  Expression Expression
}

type InputParameterType struct {}

type ExpressionToolInputs struct {}

type ExpressionToolOutputs struct {}

type ExpressionToolOutputParameter struct {
  ID string
  Label string
  Streamable bool
  Doc string

  SecondaryFiles []Expression
  Format ParameterFormat
  OutputBinding CommandOutputBinding
  Type OutputParameterType
}

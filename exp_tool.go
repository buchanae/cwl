package cwl

type ExpressionTool struct {
	Version string `cwl:"cwlVersion"`

	ID    string
	Label string
	Doc   string

	Hints        []Hint
	Requirements []Requirement

	//Inputs ExpressionToolInputs
	//Outputs ExpressionToolOutputs

	Expression Expression
}

type ExpressionToolOutput struct {
	ID         string
	Label      string
	Doc        string
	Streamable bool

	SecondaryFiles []Expression
	Format         Format
	//OutputBinding CommandOutputBinding
	Type []Type
}

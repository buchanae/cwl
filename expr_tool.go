package cwl

type ExpressionTool struct {
	CWLVersion string `json:"cwlVersion,omitempty"`
	Class      string `json:"class,omitempty"`

	ID    string `json:"id,omitempty"`
	Label string `json:"label,omitempty"`
	Doc   string `json:"doc,omitempty"`

	Hints        []Requirement `json:"hints,omitempty"`
	Requirements []Requirement `json:"requirements,omitempty"`

	Inputs  []CommandInput  `json:"inputs,omitempty"`
	Outputs []CommandOutput `json:"outputs,omitempty"`

	Expression Expression `json:"expression,omitempty"`
}

/*
type ExpressionToolInput struct {
	ID    string `json:"id,omitempty"`
	Label string `json:"label,omitempty"`
	Doc        string `json:"doc,omitempty"`
	Streamable bool   `json:"streamable,omitempty"`
	Default        Any                 `json:"default,omitempty"`

	Type           []InputType         `json:"type,omitempty"`

	SecondaryFiles []Expression        `json:"secondaryFiles,omitempty"`
	Format         []Expression        `json:"format,omitempty"`

	InputBinding   *CommandLineBinding `json:"inputBinding,omitempty"`
}

type ExpressionToolOutput struct {
	ID         string `json:"id,omitempty"`
	Label      string `json:"label,omitempty"`
	Doc        string `json:"doc,omitempty"`
	Streamable bool   `json:"streamable,omitempty"`

	Type []OutputType `json:"type,omitempty"`

	SecondaryFiles []Expression `json:"secondaryFiles,omitempty"`
	Format         []Expression `json:"format,omitempty"`

	OutputBinding *CommandOutputBinding `json:"outputBinding,omitempty"`
}
*/

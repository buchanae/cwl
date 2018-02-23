package cwl

type Workflow struct {
	Version string `cwl:"cwlVersion"`
	ID      string
	Label   string
	Doc     string

	Hints        []Hint
	Requirements []Requirement

	Inputs  []WorkflowInput
	Outputs []WorkflowOutput
	Steps   []Step
}
func (Workflow) doctype() {}

type WorkflowInput struct {
	ID         string
	Label      string
  // TODO ensure that an array of strings can be loaded
	Doc        string
	Streamable bool

	SecondaryFiles []Expression
	Format         []Expression
	InputBinding CommandLineBinding
	Default      Any
	Type         []Type
}

type WorkflowOutput struct {
	ID         string
	Label      string
	Doc        string
	Streamable bool
	LinkMerge      LinkMergeMethod

	Type           []Type
	SecondaryFiles []Expression
	Format         []Expression

	OutputBinding CommandOutputBinding
	OutputSource  []string
}

type Step struct {
	ID    string
	Label string
	Doc   string

	Hints        []Hint
	Requirements []Requirement

	In  []StepInput
	Out []StepOutput

  // TODO can be a file reference. need DocumentReference type.
	Run Document

	Scatter []string
	ScatterMethod ScatterMethod
}

type StepInput struct {
	ID        string
	Source    []string
	LinkMerge LinkMergeMethod
	Default   Any
	ValueFrom Expression
}

type StepOutput struct {
	ID string
}

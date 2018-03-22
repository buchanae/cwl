package cwl

type Workflow struct {
	CWLVersion string `json:"cwlVersion,omitempty"`
	ID         string `json:"id,omitempty"`
	Label      string `json:"label,omitempty"`
	Doc        string `json:"doc,omitempty"`

	Hints        []Requirement `json:"hints,omitempty"`
	Requirements []Requirement `json:"requirements,omitempty"`

	Inputs  []WorkflowInput  `json:"inputs,omitempty"`
	Outputs []WorkflowOutput `json:"outputs,omitempty"`
	Steps   []Step           `json:"steps,omitempty"`
}

// TODO exactly the same and CommandInput?
type WorkflowInput struct {
	ID         string `json:"id,omitempty"`
	Label      string `json:"label,omitempty"`
	Doc        string `json:"doc,omitempty"`
	Streamable bool   `json:"streamable,omitempty"`
	Default        Value               `json:"default,omitempty"`

	Type           []InputType         `json:"type,omitempty"`

	SecondaryFiles []Expression        `json:"secondaryFiles,omitempty"`
	Format         []Expression        `json:"format,omitempty"`

	InputBinding   *CommandLineBinding `json:"inputBinding,omitempty"`
}

type WorkflowOutput struct {
	ID         string          `json:"id,omitempty"`
	Label      string          `json:"label,omitempty"`
	Doc        string          `json:"doc,omitempty"`
	Streamable bool            `json:"streamable,omitempty"`
	LinkMerge  LinkMergeMethod `json:"linkMerge,omitempty"`

	Type           []OutputType `json:"type,omitempty"`
	SecondaryFiles []Expression `json:"secondaryFiles,omitempty"`
	Format         []Expression `json:"format,omitempty"`

	OutputBinding *CommandOutputBinding `json:"outputBinding,omitempty"`
	OutputSource  []string              `json:"outputSource,omitempty"`
}

type Step struct {
	ID    string `json:"id,omitempty"`
	Label string `json:"label,omitempty"`
	Doc   string `json:"doc,omitempty"`

	Hints        []Requirement `json:"hints,omitempty"`
	Requirements []Requirement `json:"requirements,omitempty"`

	In  []StepInput  `json:"in,omitempty"`
	Out []StepOutput `json:"out,omitempty"`

	Run Document `json:"run,omitempty"`

	Scatter       []string      `json:"scatter,omitempty"`
	ScatterMethod ScatterMethod `json:"scatterMethod,omitempty"`
}

type StepInput struct {
	ID        string          `json:"id,omitempty"`
	Source    []string        `json:"source,omitempty"`
	LinkMerge LinkMergeMethod `json:"linkMerge,omitempty"`
	Default   Value           `json:"default,omitempty"`
	ValueFrom Expression      `json:"valueFrom,omitempty"`
}

type StepOutput struct {
	ID string `json:"id,omitempty"`
}

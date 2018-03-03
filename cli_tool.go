package cwl

type CommandLineTool struct {
	CWLVersion string `json:"cwlVersion,omitempty"`
	Class      string `json:"class,omitempty"`
	ID         string `json:"id,omitempty"`
	Label      string `json:"label,omitempty"`
	Doc        string `json:"doc,omitempty"`

	Hints        []Hint        `json:"hints,omitempty"`
	Requirements []Requirement `json:"requirements,omitempty"`

	Inputs  []CommandInput  `json:"inputs,omitempty"`
	Outputs []CommandOutput `json:"outputs,omitempty"`

	BaseCommand []string              `json:"baseCommand,omitempty"`
	Arguments   []*CommandLineBinding `json:"arguments,omitempty"`

	Stdin  Expression `json:"stdin,omitempty"`
	Stderr Expression `json:"stderr,omitempty"`
	Stdout Expression `json:"stdout,omitempty"`

	SuccessCodes       []int `json:"successCodes,omitempty"`
	TemporaryFailCodes []int `json:",omitempty"`
	PermanentFailCodes []int `json:",omitempty"`
}

type CommandInput struct {
	ID         string `json:"id,omitempty"`
	Label      string `json:"label,omitempty"`
	Doc        string `json:"doc,omitempty"`
	Streamable bool   `json:"streamable,omitempty"`
	Default    Any    `json:"default,omitempty"`

	Type []InputType `json:"type,omitempty"`

	SecondaryFiles []Expression `json:"secondaryFiles,omitempty"`
	Format         []Expression `json:"format,omitempty"`

	InputBinding *CommandLineBinding `json:"inputBinding,omitempty"`
}

type CommandOutput struct {
	ID         string `json:"id,omitempty"`
	Label      string `json:"label,omitempty"`
	Doc        string `json:"doc,omitempty"`
	Streamable bool   `json:"streamable,omitempty"`

	Type []OutputType `json:"type,omitempty"`

	SecondaryFiles []Expression `json:"secondaryFiles,omitempty"`
	Format         []Expression `json:"format,omitempty"`

	OutputBinding *CommandOutputBinding `json:"outputBinding,omitempty"`
}

type CommandLineBinding struct {
	LoadContents  bool       `json:"loadContents,omitempty"`
	Position      int        `json:"position,omitempty"`
	Prefix        string     `json:"prefix,omitempty"`
	ItemSeparator string     `json:"itemSeparator,omitempty"`
	ValueFrom     Expression `json:"valueFrom,omitempty"`

	separate      bool
	separateSet   bool
	shellQuote    bool
	shellQuoteSet bool
}

func (c *CommandLineBinding) Separate() bool {
	if !c.separateSet {
		return true
	}
	return c.separate
}

func (c *CommandLineBinding) SetSeparate(b bool) {
	c.separate = b
	c.separateSet = true
}

func (c *CommandLineBinding) ShellQuote() bool {
	if !c.shellQuoteSet {
		return true
	}
	return c.shellQuote
}

func (c *CommandLineBinding) SetShellQuote(b bool) {
	c.shellQuote = b
	c.shellQuoteSet = true
}

type CommandOutputBinding struct {
	Glob         []Expression `json:"glob,omitempty"`
	LoadContents bool         `json:"loadContents,omitempty"`
	OutputEval   Expression   `json:"outputEval,omitempty"`
}

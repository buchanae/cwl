package cwl

type CommandLineTool struct {
	Version string `cwl:"cwlVersion"`
	ID      string
	Label   string
	Doc     string

	Hints        []Hint
	Requirements []Requirement

	Inputs  []CommandInput
	Outputs []CommandOutput

	BaseCommand []string
	Arguments   []CommandLineBinding

	Stdin  Expression
	Stderr Expression
	Stdout Expression

	SuccessCodes       []int
	TemporaryFailCodes []int
	PermanentFailCodes []int
}

type CommandLineBinding struct {
	LoadContents  bool
	Position      int
	Prefix        string
	ItemSeparator string
	ValueFrom     Expression

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

type CommandInput struct {
	ID         string
	Label      string
	Doc        string
	Streamable bool
	Default    Any

	Type []InputType

	SecondaryFiles []Expression
	Format         []Expression

	InputBinding CommandLineBinding
}

type CommandOutput struct {
	ID         string
	Label      string
	Doc        string
	Streamable bool

	Type []OutputType

	SecondaryFiles []Expression
	Format         []Expression

	OutputBinding CommandOutputBinding
}

type CommandOutputBinding struct {
	Glob         []Expression
	LoadContents bool
	OutputEval   Expression
}

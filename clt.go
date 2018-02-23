package cwl

type Document interface {
	doctype()
}

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

func (CommandLineTool) doctype() {}

type CommandInput struct {
	ID         string
	Label      string
	Doc        string
	Streamable bool
	Default    Any

	Type []Type

	SecondaryFiles []Expression
	Format         []Expression

	InputBinding CommandLineBinding
}

type CommandLineBinding struct {
	LoadContents  bool
	Position      int
	Prefix        string
	Separate      bool
	ItemSeparator string
	ValueFrom     Expression
	ShellQuote    bool
}

type CommandOutput struct {
	ID         string
	Label      string
	Doc        string
	Streamable bool

	Type []Type

	SecondaryFiles []Expression
	Format         []Expression

	OutputBinding CommandOutputBinding
}

type CommandOutputBinding struct {
	Glob         []Expression
	LoadContents bool
	OutputEval   Expression
}

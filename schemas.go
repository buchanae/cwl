package cwl

type InputRecord struct {
	Label  string
	Fields []InputRecordField
}

type InputRecordField struct {
	Name         string
	Doc          string
	Label        string
	Type         []InputType
	InputBinding CommandLineBinding
}

type InputEnum struct {
	Label        string
	Symbols      []string
	InputBinding CommandLineBinding
}

type InputArray struct {
	Label        string
	Items        []InputType
	InputBinding CommandLineBinding
}

type OutputRecord struct {
	Label  string
	Fields []OutputRecordField
}

type OutputRecordField struct {
	Name          string
	Doc           string
	Type          []OutputType
	OutputBinding CommandOutputBinding
}

type OutputEnum struct {
	Label         string
	Symbols       []string
	OutputBinding CommandOutputBinding
}

type OutputArray struct {
	Label         string
	Items         []OutputType
	OutputBinding CommandOutputBinding
}

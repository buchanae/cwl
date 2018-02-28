package cwl

type InputRecord struct {
	Label  string       `json:"label,omitempty"`
	Fields []InputField `json:"fields,omitempty"`
}

type InputField struct {
	Name         string              `json:"name,omitempty"`
	Doc          string              `json:"doc,omitempty"`
	Label        string              `json:"label,omitempty"`
	Type         []InputType         `json:"type,omitempty"`
	InputBinding *CommandLineBinding `json:"inputBinding,omitempty"`
}

type InputEnum struct {
	Label        string              `json:"label,omitempty"`
	Symbols      []string            `json:"symbols,omitempty"`
	InputBinding *CommandLineBinding `json:"inputBinding,omitempty"`
}

type InputArray struct {
	Label        string              `json:"label,omitempty"`
	Items        []InputType         `json:"items,omitempty"`
	InputBinding *CommandLineBinding `json:"inputBinding,omitempty"`
}

type OutputRecord struct {
	Label  string        `json:"label,omitempty"`
	Fields []OutputField `json:"fields,omitempty"`
}

type OutputField struct {
	Name          string                `json:"name,omitempty"`
	Doc           string                `json:"doc,omitempty"`
	Type          []OutputType          `json:"type,omitempty"`
	OutputBinding *CommandOutputBinding `json:"outputBinding,omitempty"`
}

type OutputEnum struct {
	Label         string                `json:"label,omitempty"`
	Symbols       []string              `json:"symbols,omitempty"`
	OutputBinding *CommandOutputBinding `json:"outputBinding,omitempty"`
}

type OutputArray struct {
	Label         string                `json:"label,omitempty"`
	Items         []OutputType          `json:"items,omitempty"`
	OutputBinding *CommandOutputBinding `json:"outputBinding,omitempty"`
}

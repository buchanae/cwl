package cwl

type OutputRecordSchema struct {
	Type   string
	Fields []OutputRecordField
	Label  string
}

type OutputRecordFieldType struct {
}

type OutputRecordField struct {
	Name string
	Type OutputRecordFieldType
	Doc  string
	OutputBinding CommandOutputBinding
}

type OutputEnumSchema struct {
	Symbols []string
	Type    string
	Label   string
	OutputBinding CommandOutputBinding
}

type OutputArraySchemaItems struct {
	Type
}

type OutputArraySchema struct {
	Items OutputArraySchemaItems
	Type  string
	Label string
	OutputBinding CommandOutputBinding
}

type InputSchema struct {
}

type InputRecordSchema struct {
	Type   string
	Fields []InputRecordField
	Label  string
}

type InputRecordFieldType struct {
}

type InputRecordField struct {
	Name string
	Type InputRecordFieldType
	Doc  string
	InputBinding CommandLineBinding
	Label string
}

type InputEnumSchema struct {
	Symbols []string
	Type    string
	Label   string
	InputBinding CommandLineBinding
}

type InputArraySchemaItems struct{}

type InputArraySchema struct {
	Items InputArraySchemaItems
	Type  string
	Label string
	InputBinding CommandLineBinding
}

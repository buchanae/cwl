package cwl

/*
cwlVersion: v1.0
class: CommandLineTool
inputs:
  dependent_parameters:
    type:
      type: record
      name: dependent_parameters
      fields:
        itemA:
          type: string
          inputBinding:
            prefix: -A
        itemB:
          type: string
          inputBinding:
            prefix: -B
  exclusive_parameters:
    type:
      - type: record
        name: itemC
        fields:
          itemC:
            type: string
            inputBinding:
              prefix: -C
      - type: record
        name: itemD
        fields:
          itemD:
            type: string
            inputBinding:
              prefix: -D
outputs: []
baseCommand: echo


dependent_parameters:
  itemA: one
exclusive_parameters:
  itemC: three

$ cwl-runner record.cwl record-job1.yml
Workflow error:
  Error validating input record, could not validate field `dependent_parameters` because
  missing required field `itemB`

dependent_parameters:
  itemA: one
  itemB: two
exclusive_parameters:
  itemC: three
  itemD: four

$ cwl-runner record.cwl record-job2.yml
[job 140566927111376] /home/example$ echo -A one -B two -C three
-A one -B two -C three
Final process status is success
{}


dependent_parameters:
  itemA: one
  itemB: two
exclusive_parameters:
  itemD: four

$ cwl-runner record.cwl record-job3.yml
[job 140606932172880] /home/example$ echo -A one -B two -D four
-A one -B two -D four
Final process status is success
{}
*/
type InputSchema struct {
}

type InputRecordSchema struct {
	Type   string
	Fields []InputRecordField
	Label  string
}

type InputRecordField struct {
	Name         string
	Type         []Type
	Doc          string
	InputBinding CommandLineBinding
	Label        string
}

type InputEnumSchema struct {
	Symbols      []string
	Type         string
	Label        string
	InputBinding CommandLineBinding
}

type InputArraySchemaItems struct{}

type InputArraySchema struct {
	Items        InputArraySchemaItems
	Type         string
	Label        string
	InputBinding CommandLineBinding
}

type OutputRecordSchema struct {
	Type   string
	Fields []OutputRecordField
	Label  string
}

type OutputRecordField struct {
	Name          string
	Type          []Type
	Doc           string
	OutputBinding CommandOutputBinding
}

type OutputEnumSchema struct {
	Symbols       []string
	Type          string
	Label         string
	OutputBinding CommandOutputBinding
}

type OutputArraySchemaItems struct {
	Type
}

type OutputArraySchema struct {
	Items         OutputArraySchemaItems
	Type          string
	Label         string
	OutputBinding CommandOutputBinding
}

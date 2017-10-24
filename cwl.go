package cwl

import (
  "github.com/ghodss/yaml"
  "encoding/json"
)

func Parse(raw string) interface{} {
  r, _ := ParseWorkflow(raw)
  return r
}

func ParseWorkflow(raw string) (*Workflow, error) {
  wf := new(Workflow)
  err := parse(raw, wf)
  return wf, err
}

func ParseTool(raw string) (*Tool, error) {
  tool := new(Tool)
  err := parse(raw, tool)
  return tool, err
}

func parse(source string, doc interface{}) error {
  jsonb, yerr := yaml.YAMLToJSON([]byte(source))
  if yerr != nil {
    return yerr
  }
  err := json.Unmarshal(jsonb, &doc)
	if err != nil {
		return err
	}
	return nil
}

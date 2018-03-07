package cwl

import (
	"encoding/json"
)

// A bunch of tedious wrappers for fields like "class" and "type"
// so that they marhshal to JSON/YAML correctly.

func (i File) MarshalJSON() ([]byte, error) {
	type Wrap File
	return json.Marshal(struct {
		Type string `json:"type"`
		Wrap
	}{"File", Wrap(i)})
}
func (i Directory) MarshalJSON() ([]byte, error) {
	type Wrap Directory
	return json.Marshal(struct {
		Type string `json:"type"`
		Wrap
	}{"Directory", Wrap(i)})
}

func (i InputArray) MarshalJSON() ([]byte, error) {
	type Wrap InputArray
	return json.Marshal(struct {
		Type string `json:"type"`
		Wrap
	}{"array", Wrap(i)})
}

func (i OutputArray) MarshalJSON() ([]byte, error) {
	type Wrap OutputArray
	return json.Marshal(struct {
		Type string `json:"type"`
		Wrap
	}{"array", Wrap(i)})
}
func (i InputRecord) MarshalJSON() ([]byte, error) {
	type Wrap InputRecord
	return json.Marshal(struct {
		Type string `json:"type"`
		Wrap
	}{"record", Wrap(i)})
}

func (i OutputRecord) MarshalJSON() ([]byte, error) {
	type Wrap OutputRecord
	return json.Marshal(struct {
		Type string `json:"type"`
		Wrap
	}{"record", Wrap(i)})
}
func (x Workflow) MarshalJSON() ([]byte, error) {
	type Wrap Workflow
	return json.Marshal(struct {
		Class string `json:"class"`
		Wrap
	}{"Workflow", Wrap(x)})
}

func (x Tool) MarshalJSON() ([]byte, error) {
	type Wrap Tool
	return json.Marshal(struct {
		Class string `json:"class"`
		Wrap
	}{"CommandLineTool", Wrap(x)})
}

func (x DockerRequirement) MarshalJSON() ([]byte, error) {
	type Wrap DockerRequirement
	return json.Marshal(struct {
		Class string `json:"class"`
		Wrap
	}{"DockerRequirement", Wrap(x)})
}
func (x ResourceRequirement) MarshalJSON() ([]byte, error) {
	type Wrap ResourceRequirement
	return json.Marshal(struct {
		Class string `json:"class"`
		Wrap
	}{"ResourceRequirement", Wrap(x)})
}
func (x EnvVarRequirement) MarshalJSON() ([]byte, error) {
	type Wrap EnvVarRequirement
	return json.Marshal(struct {
		Class string `json:"class"`
		Wrap
	}{"EnvVarRequirement", Wrap(x)})
}
func (x SchemaDefRequirement) MarshalJSON() ([]byte, error) {
	type Wrap SchemaDefRequirement
	return json.Marshal(struct {
		Class string `json:"class"`
		Wrap
	}{"SchemaDefRequirement", Wrap(x)})
}
func (x ShellCommandRequirement) MarshalJSON() ([]byte, error) {
	type Wrap ShellCommandRequirement
	return json.Marshal(struct {
		Class string `json:"class"`
		Wrap
	}{"ShellCommandRequirement", Wrap(x)})
}
func (x InlineJavascriptRequirement) MarshalJSON() ([]byte, error) {
	type Wrap InlineJavascriptRequirement
	return json.Marshal(struct {
		Class string `json:"class"`
		Wrap
	}{"InlineJavascriptRequirement", Wrap(x)})
}
func (x SoftwareRequirement) MarshalJSON() ([]byte, error) {
	type Wrap SoftwareRequirement
	return json.Marshal(struct {
		Class string `json:"class"`
		Wrap
	}{"SoftwareRequirement", Wrap(x)})
}
func (x InitialWorkDirRequirement) MarshalJSON() ([]byte, error) {
	type Wrap InitialWorkDirRequirement
	return json.Marshal(struct {
		Class string `json:"class"`
		Wrap
	}{"InitialWorkDirRequirement", Wrap(x)})
}
func (x SubworkflowFeatureRequirement) MarshalJSON() ([]byte, error) {
	type Wrap SubworkflowFeatureRequirement
	return json.Marshal(struct {
		Class string `json:"class"`
		Wrap
	}{"SubworkflowFeatureRequirement", Wrap(x)})
}
func (x ScatterFeatureRequirement) MarshalJSON() ([]byte, error) {
	type Wrap ScatterFeatureRequirement
	return json.Marshal(struct {
		Class string `json:"class"`
		Wrap
	}{"ScatterFeatureRequirement", Wrap(x)})
}
func (x MultipleInputFeatureRequirement) MarshalJSON() ([]byte, error) {
	type Wrap MultipleInputFeatureRequirement
	return json.Marshal(struct {
		Class string `json:"class"`
		Wrap
	}{"MultipleInputFeatureRequirement", Wrap(x)})
}
func (x StepInputExpressionRequirement) MarshalJSON() ([]byte, error) {
	type Wrap StepInputExpressionRequirement
	return json.Marshal(struct {
		Class string `json:"class"`
		Wrap
	}{"StepInputExpressionRequirement", Wrap(x)})
}

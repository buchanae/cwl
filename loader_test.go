package cwl

import (
	"github.com/kr/pretty"
	"reflect"
	"testing"
)

func TestLoadSimpleFile(t *testing.T) {
	doc, err := LoadFile("./examples/1st-tool.yml")
	if err != nil {
		t.Fatal(err)
	}

	c := doc.(*CommandLineTool)
	pretty.Println(c)
	e := &CommandLineTool{
		Version:     "v1.0",
		BaseCommand: []string{"echo"},
		Inputs: []CommandInput{
			{
				ID:   "message",
				Type: []InputType{String{}},
				InputBinding: CommandLineBinding{
					Position: 1,
				},
			},
		},
	}

	if !reflect.DeepEqual(c, e) {
		t.Error("different docs")
		diff := pretty.Diff(c, e)
		for _, d := range diff {
			t.Log(d)
		}
	}
}

func TestLoadSimpleFile2(t *testing.T) {
	doc, err := LoadFile("./examples/tar.cwl")
	if err != nil {
		t.Fatal(err)
	}

	c := doc.(*CommandLineTool)
	pretty.Println(c)
	e := &CommandLineTool{
		Version:     "v1.0",
		BaseCommand: []string{"tar", "xf"},
		Inputs: []CommandInput{
			{
				ID:   "tarfile",
				Type: []InputType{FileType{}},
				InputBinding: CommandLineBinding{
					Position: 1,
				},
			},
		},
		Outputs: []CommandOutput{
			{
				ID:   "example_out",
				Type: []OutputType{FileType{}},
				OutputBinding: CommandOutputBinding{
					Glob: []Expression{"hello.txt"},
				},
			},
		},
	}

	if !reflect.DeepEqual(c, e) {
		t.Error("different docs")
		diff := pretty.Diff(c, e)
		for _, d := range diff {
			t.Log(d)
		}
	}
}

func TestStableMapOrder(t *testing.T) {
	for i := 0; i < 20; i++ {
		t.Run("", TestLoadCltAll)
	}
}

func TestLoadMC3Wf(t *testing.T) {
	_, err := LoadFile("./examples/mc3-annotate.cwl")
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadSimpleWf(t *testing.T) {
	d, err := LoadFile("./examples/1st-workflow.yml")
	if err != nil {
		t.Fatal(err)
	}
	e := &Workflow{
		Version:      "v1.0",
		ID:           "",
		Label:        "",
		Doc:          "",
		Hints:        nil,
		Requirements: nil,
		Inputs: []WorkflowInput{
			{
				ID:             "inp",
				Label:          "",
				Doc:            "",
				Streamable:     false,
				SecondaryFiles: nil,
				Format:         nil,
				InputBinding:   CommandLineBinding{},
				Default:        nil,
				Type: []InputType{
					FileType{},
				},
			},
			{
				ID:             "ex",
				Label:          "",
				Doc:            "",
				Streamable:     false,
				SecondaryFiles: nil,
				Format:         nil,
				InputBinding:   CommandLineBinding{},
				Default:        nil,
				Type: []InputType{
					String{},
				},
			},
			{
				ID:             "foo",
				Label:          "",
				Doc:            "doc1\ndoc2",
				Streamable:     false,
				SecondaryFiles: []Expression{".bai"},
				Format:         []Expression{"fmt"},
				InputBinding:   CommandLineBinding{},
				Default:        nil,
				Type:           nil,
			},
			{
				ID:             "bar",
				Label:          "",
				Doc:            "docstring",
				Streamable:     false,
				SecondaryFiles: []Expression{".fai", ".bai"},
				Format:         []Expression{"fm1", "fm2"},
				InputBinding:   CommandLineBinding{},
				Default:        nil,
				Type:           nil,
			},
		},
		Outputs: []WorkflowOutput{
			{
				ID:         "other",
				Label:      "",
				Doc:        "",
				Streamable: false,
				LinkMerge:  0,
				Type: []OutputType{
					OutputArray{
						Items: []OutputType{
							FileType{},
						},
						OutputBinding: CommandOutputBinding{},
					},
				},
				SecondaryFiles: nil,
				Format:         nil,
				OutputBinding:  CommandOutputBinding{},
				OutputSource:   nil,
			},
			{
				ID:         "classout",
				Label:      "",
				Doc:        "",
				Streamable: false,
				LinkMerge:  0,
				Type: []OutputType{
					FileType{},
				},
				SecondaryFiles: nil,
				Format:         nil,
				OutputBinding:  CommandOutputBinding{},
				OutputSource:   []string{"compile/classfile"},
			},
		},
		Steps: []Step{
			{
				ID:           "subwf",
				Label:        "",
				Doc:          "",
				Hints:        nil,
				Requirements: nil,
				In:           []StepInput{},
				Out: []StepOutput{
					{ID: "one"},
				},
				Run: &CommandLineTool{
					Version: "",
					ID:      "",
					Label:   "",
					Doc:     "",
					Hints:   nil,
					Requirements: []Requirement{
						ShellCommandRequirement{},
					},
					Inputs: nil,
					Outputs: []CommandOutput{
						{
							ID:         "one",
							Label:      "",
							Doc:        "doc1\ndoc2",
							Streamable: false,
							Type: []OutputType{
								FileType{},
							},
							SecondaryFiles: []Expression{".foo"},
							Format:         []Expression{"fmt"},
							OutputBinding: CommandOutputBinding{
								Glob:         []Expression{"*.glob"},
								LoadContents: false,
								OutputEval:   "",
							},
						},
						{
							ID:         "arrouttest",
							Label:      "",
							Doc:        "docstring",
							Streamable: false,
							Type: []OutputType{
								OutputArray{
									Items: []OutputType{
										FileType{},
									},
									OutputBinding: CommandOutputBinding{},
								},
							},
							SecondaryFiles: []Expression{".fai", ".bai"},
							Format:         []Expression{"fm1", "fm2"},
							OutputBinding: CommandOutputBinding{
								Glob:         []Expression{"*.glob1", "*.glob2"},
								LoadContents: false,
								OutputEval:   "",
							},
						},
					},
					BaseCommand: nil,
					Arguments: []CommandLineBinding{
						{
							LoadContents:  false,
							Position:      0,
							Prefix:        "",
							Separate:      true,
							ItemSeparator: "",
							ValueFrom:     "date\ntar cf hello.tar Hello.java\ndate\n",
							ShellQuote:    false,
						},
					},
					Stdin:              "",
					Stderr:             "",
					Stdout:             "",
					SuccessCodes:       nil,
					TemporaryFailCodes: nil,
					PermanentFailCodes: nil,
				},
				Scatter:       nil,
				ScatterMethod: 0,
			},
			{
				ID:           "auntar",
				Label:        "",
				Doc:          "",
				Hints:        nil,
				Requirements: nil,
				In: []StepInput{
					{
						ID:        "tarfile",
						Source:    []string{"inp"},
						LinkMerge: 0,
						Default:   nil,
						ValueFrom: "",
					},
					{
						ID:        "other",
						Source:    []string{"ex"},
						LinkMerge: 0,
						Default:   nil,
						ValueFrom: "",
					},
					{
						ID:        "extractfile",
						Source:    []string{"ex"},
						LinkMerge: 0,
						Default:   nil,
						ValueFrom: "",
					},
				},
				Out: []StepOutput{
					{ID: "example_out"},
				},
				Run:           DocumentRef{URL: "tar-param.cwl"},
				Scatter:       []string{"tarfile"},
				ScatterMethod: 0,
			},
			{
				ID:           "untar",
				Label:        "",
				Doc:          "",
				Hints:        nil,
				Requirements: nil,
				In: []StepInput{
					{
						ID:        "tarfile",
						Source:    []string{"inp"},
						LinkMerge: 0,
						Default:   nil,
						ValueFrom: "",
					},
					{
						ID:        "extractfile",
						Source:    []string{"ex"},
						LinkMerge: 0,
						Default:   nil,
						ValueFrom: "",
					},
				},
				Out: []StepOutput{
					{ID: "example_out"},
				},
				Run:           DocumentRef{URL: "tar-param.cwl"},
				Scatter:       nil,
				ScatterMethod: 0,
			},
			{
				ID:           "compile",
				Label:        "",
				Doc:          "",
				Hints:        nil,
				Requirements: nil,
				In: []StepInput{
					{
						ID:        "src",
						Source:    []string{"untar/example_out"},
						LinkMerge: 0,
						Default:   nil,
						ValueFrom: "",
					},
				},
				Out: []StepOutput{
					{ID: "classfile"},
				},
				Run:           DocumentRef{URL: "arguments.cwl"},
				Scatter:       nil,
				ScatterMethod: 0,
			},
		},
	}
	pretty.Println(d)

	if !reflect.DeepEqual(d, e) {
		t.Error("different docs")
		diff := pretty.Diff(d, e)
		for _, di := range diff {
			t.Log(di)
		}
	}
}

func TestLoadCltAll(t *testing.T) {
	doc, err := LoadFile("./examples/clt-all.cwl")
	if err != nil {
		t.Fatal(err)
	}

	c := doc.(*CommandLineTool)
	pretty.Println(c)
	e := &CommandLineTool{
		Version:     "v1.0",
		Label:       `Example trivial wrapper for Java 7 compiler`,
		Doc:         "Example doc",
		BaseCommand: []string{"echo", "foo"},
		Hints: []Hint{
			DockerRequirement{
				Pull: "java:7-jdk",
			},
			DockerRequirement{
				Load: "loadjava:7-jdk",
			},
		},
		Arguments: []CommandLineBinding{
			{
				ValueFrom: Expression("-d"),
			},
			{
				ValueFrom: Expression("$(runtime.outdir)"),
			},
		},
		Stdout: Expression("output.txt"),
		Stderr: Expression("error.txt"),
		Inputs: []CommandInput{
			{
				ID:   "tarfile",
				Type: []InputType{FileType{}},
				InputBinding: CommandLineBinding{
					Position: 1,
				},
			},
			{
				ID:   "extractfile",
				Type: []InputType{String{}},
				InputBinding: CommandLineBinding{
					Position: 2,
				},
			},
			{
				ID:   "nullablefile",
				Type: []InputType{Null{}, String{}},
				InputBinding: CommandLineBinding{
					Position: 2,
				},
			},
			{
				ID: "list",
				Type: []InputType{
					InputArray{
						Items: []InputType{
							String{},
						},
					},
				},
				InputBinding: CommandLineBinding{
					Position:      3,
					ItemSeparator: ",",
					Separate:      true,
					Prefix:        "-A",
				},
			},
			{
				ID: "list2",
				Type: []InputType{
					InputArray{
						Items: []InputType{
							String{},
						},
					},
				},
			},
			{
				ID:   "optional_file",
				Type: []InputType{FileType{}, Null{}},
			},
			{
				ID:   "flag",
				Type: []InputType{Boolean{}},
			},
			{
				ID:   "num",
				Type: []InputType{Int{}},
			},
		},
		Outputs: []CommandOutput{
			{
				ID:   "output1",
				Type: []OutputType{Stdout{}},
			},
			{
				ID:   "error1",
				Type: []OutputType{Stderr{}},
			},
			{
				ID:   "example_out",
				Type: []OutputType{FileType{}},
				OutputBinding: CommandOutputBinding{
					Glob: []Expression{"$(inputs.extractfile)"},
				},
			},
			{
				ID: "arrayoutput",
				Type: []OutputType{
					OutputArray{
						Items: []OutputType{
							String{},
						},
					},
				},
			},
			{
				ID: "arrayoutput2",
				Type: []OutputType{
					OutputArray{
						Items: []OutputType{
							String{},
						},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(c, e) {
		t.Error("different docs")
		diff := pretty.Diff(c, e)
		for _, d := range diff {
			t.Log(d)
		}
	}
}

func TestUnknownClass(t *testing.T) {
	_, err := Load([]byte(`class: Foo`))
	if err == nil {
		t.Error("expected error")
	}
	t.Log(err)
}

func TestInvalidYAML(t *testing.T) {
	_, err := Load([]byte(`:f`))
	if err == nil {
		t.Error("expected error")
	}
	t.Log(err)
}

func TestEmptyYAML(t *testing.T) {
	_, err := Load([]byte{})
	if err == nil {
		t.Error("expected error")
	}
	t.Log(err)
}

func TestYAMLList(t *testing.T) {
	_, err := Load([]byte(`- foo\n-foo`))
	if err == nil {
		t.Error("expected error")
	}
	t.Log(err)
}

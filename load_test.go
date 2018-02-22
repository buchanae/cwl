package cwl

import (
	"github.com/go-test/deep"
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
				Type: []Type{String},
				InputBinding: CommandLineBinding{
					Position: 1,
				},
			},
		},
	}

	if !reflect.DeepEqual(c, e) {
		t.Error("different docs")
		diff := deep.Equal(c, e)
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
				Type: []Type{FileType},
				InputBinding: CommandLineBinding{
					Position: 1,
				},
			},
		},
		Outputs: []CommandOutput{
			{
				ID:   "example_out",
				Type: []Type{FileType},
				OutputBinding: CommandOutputBinding{
					Glob: []Expression{"hello.txt"},
				},
			},
		},
	}

	if !reflect.DeepEqual(c, e) {
		t.Error("different docs")
		diff := deep.Equal(c, e)
		for _, d := range diff {
			t.Log(d)
		}
	}
}

func TestStableMapOrder(t *testing.T) {
	for i := 0; i < 20; i++ {
		t.Run("", TestLoadSimpleFile3)
	}
}

func TestLoadSimpleFile3(t *testing.T) {
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
				Type: []Type{FileType},
				InputBinding: CommandLineBinding{
					Position: 1,
				},
			},
			{
				ID:   "extractfile",
				Type: []Type{String},
				InputBinding: CommandLineBinding{
					Position: 2,
				},
			},
			{
				ID:   "nullablefile",
				Type: []Type{Null, String},
				InputBinding: CommandLineBinding{
					Position: 2,
				},
			},
			{
				ID:   "list",
				Type: []Type{ArrayType{String}},
				InputBinding: CommandLineBinding{
					Position:      3,
					ItemSeparator: ",",
					Separate:      true,
					Prefix:        "-A",
				},
			},
			{
				ID:   "list2",
				Type: []Type{ArrayType{String}},
			},
			{
				ID:   "optional_file",
				Type: []Type{FileType, Null},
			},
			{
				ID:   "flag",
				Type: []Type{Boolean},
			},
			{
				ID:   "num",
				Type: []Type{Int},
			},
		},
		Outputs: []CommandOutput{
			{
				ID:   "output1",
				Type: []Type{Stdout},
			},
			{
				ID:   "error1",
				Type: []Type{Stderr},
			},
			{
				ID:   "example_out",
				Type: []Type{FileType},
				OutputBinding: CommandOutputBinding{
					Glob: []Expression{"$(inputs.extractfile)"},
				},
			},
			{
				ID:   "arrayoutput",
				Type: []Type{ArrayType{FileType}},
			},
			{
				ID:   "arrayoutput2",
				Type: []Type{ArrayType{String}},
			},
		},
	}

	if !reflect.DeepEqual(c, e) {
		t.Error("different docs")
		diff := deep.Equal(c, e)
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

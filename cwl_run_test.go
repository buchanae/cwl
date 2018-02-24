package cwl

import (
	"testing"
)

func TestBuildCommand(t *testing.T) {
	// TODO load these cli arg tests from yaml docs
	clt := CommandLineTool{
		BaseCommand: []string{"echo", "hello"},
		Inputs: []CommandInput{
			{
				ID:   "test1",
				Type: []InputType{String{}},
				InputBinding: CommandLineBinding{
					Position: 3,
				},
			},
			{
				ID:   "test2",
				Type: []InputType{String{}},
				InputBinding: CommandLineBinding{
					Position: 2,
				},
			},
			{
				ID:   "test3",
				Type: []InputType{String{}, Null{}},
			},
			{
				ID:   "test4",
				Type: []InputType{String{}},
				InputBinding: CommandLineBinding{
					Position: 3,
					Prefix:   "-B",
				},
			},
			{
				ID:   "test5",
				Type: []InputType{Int{}},
				InputBinding: CommandLineBinding{
					Position: 4,
					Prefix:   "-C",
				},
			},
			{
				ID:   "test6",
				Type: []InputType{Int{}},
				InputBinding: CommandLineBinding{
					Position: 5,
					Prefix:   "-D",
					Separate: true,
				},
			},
			{
				ID:   "test7",
				Type: []InputType{Boolean{}},
				InputBinding: CommandLineBinding{
					Position: 6,
					Prefix:   "-E",
				},
			},
			{
				ID:   "test8",
				Type: []InputType{Boolean{}},
				InputBinding: CommandLineBinding{
					Position: 7,
					Prefix:   "-F",
				},
			},
			{
				ID:   "test9",
				Type: []InputType{InputArray{Items: []InputType{String{}}}},
				InputBinding: CommandLineBinding{
					Position: 8,
					Prefix:   "-G",
				},
			},
		},
	}

	err := buildCommand(clt, map[string]interface{}{
		"test1": "foo",
		"test2": "bar",
		"test4": "baz",
		"test5": "6",
		"test6": "7",
		"test7": "true",
		"test8": "false",
		"test9": []string{
			"foo", "bar",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
}

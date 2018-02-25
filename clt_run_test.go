package cwl

import (
	"fmt"
	"strings"
	"testing"
)

func TestBuildCommand(t *testing.T) {
	// TODO load these cli arg tests from yaml docs
	doc, err := LoadFile("./examples/record-clt.yml")
	if err != nil {
		t.Fatal(err)
	}
	clt := doc.(*CommandLineTool)
	debug(clt)

	args, err := buildCommand(clt, map[string]interface{}{
		"arrparam": []interface{}{"five", "six", 1},
		"arrrec": []interface{}{
			map[string]interface{}{
				"recA": "bar",
				"recB": "baz",
			},
		},
		"nil":    "foo",
		"flag":   false,
		"onflag": true,
		"zdependent_parameters": map[string]interface{}{
			"itemA": "one",
			"itemB": "two",
		},
		"exclusive_parameters": map[string]interface{}{
			"itemC": "three",
			"itemD": "four",
		},
		"extra": "foo",
	})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(strings.Join(args, " "))
}

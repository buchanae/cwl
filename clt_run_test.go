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

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(strings.Join(args, " "))
}

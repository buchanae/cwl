package cwl

import (
	"github.com/buchanae/cwl/fs"
	"testing"
)

func TestResolveFile(t *testing.T) {
	f := File{
		Location: "./examples/record.cwl",
	}

	l := fs.NewLocal()
	x, err := ResolveFile(f, l, true)
	if err != nil {
		t.Error(err)
	}

	debug(x)
}

func TestEvalSecondaryFilesPattern(t *testing.T) {

	tests := []struct {
		path    string
		pattern string
		expect  string
	}{
		{
			path:    "foo.bam",
			pattern: "^.bai",
			expect:  "foo.bai",
		},
		{
			path:    "foo.bam",
			pattern: ".bai",
			expect:  "foo.bam.bai",
		},
	}

	for _, test := range tests {
		res := EvalSecondaryFilesPattern(test.path, test.pattern)
		t.Log(test.path, test.pattern, test.expect, res)
		if res != test.expect {
			t.Error("failed match", res, test.expect)
		}
	}

}

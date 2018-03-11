package expr

import (
	"github.com/kr/pretty"
	"reflect"
	"testing"
)

func TestParseString(t *testing.T) {
	tests := []struct {
		input  string
		expect []*Part
	}{
		{
			input: "",
		},
		{
			input: "none",
			expect: []*Part{
				{Raw: "none", Start: 0, End: 4},
			},
		},
		{
			input: "$(inputs.one.path)",
			expect: []*Part{
				{
					Raw:   "$(inputs.one.path)",
					Expr:  "inputs.one.path",
					Start: 0,
					End:   18,
				},
			},
		},
		{
			input: `before $(runtime["cores"]) after`,
			expect: []*Part{
				{Raw: "before ", Start: 0, End: 7},
				{
					Raw:   `$(runtime["cores"])`,
					Expr:  `runtime["cores"]`,
					Start: 7,
					End:   26,
				},
				{Raw: " after", Start: 26, End: 32},
			},
		},
		{
			input: `before $(runtime['cores']) after`,
			expect: []*Part{
				{Raw: "before ", Start: 0, End: 7},
				{
					Raw:   `$(runtime['cores'])`,
					Expr:  `runtime['cores']`,
					Start: 7,
					End:   26,
				},
				{Raw: " after", Start: 26, End: 32},
			},
		},
		{
			input: "before $(runtime.cores[0]) after",
			expect: []*Part{
				{Raw: "before ", Start: 0, End: 7},
				{
					Raw:   `$(runtime.cores[0])`,
					Expr:  `runtime.cores[0]`,
					Start: 7,
					End:   26,
				},
				{Raw: " after", Start: 26, End: 32},
			},
		},
		{
			input: "before $(inputs.one.path) after $(two) after2",
			expect: []*Part{
				{Raw: "before ", Start: 0, End: 7},
				{
					Raw:   `$(inputs.one.path)`,
					Expr:  `inputs.one.path`,
					Start: 7,
					End:   25,
				},
				{Raw: " after ", Start: 25, End: 32},
				{
					Raw:   `$(two)`,
					Expr:  `two`,
					Start: 32,
					End:   38,
				},
				{Raw: " after2", Start: 38, End: 45},
			},
		},
		{
			input: "before $(inputs.one.path) after (two) after2",
			expect: []*Part{
				{Raw: "before ", Start: 0, End: 7},
				{
					Raw:   `$(inputs.one.path)`,
					Expr:  `inputs.one.path`,
					Start: 7,
					End:   25,
				},
				{Raw: " after (two) after2", Start: 25, End: 44},
			},
		},
		{
			input: "$()",
			expect: []*Part{
				{Raw: "$()", Expr: "", Start: 0, End: 3},
			},
		},
		{
			input: "${}",
			expect: []*Part{
				{Raw: "${}", Expr: "", Start: 0, End: 3},
			},
		},
		{
			input: "${foo bar $(bas)}",
			expect: []*Part{
				{
					Raw:   "${foo bar $(bas)}",
					Expr:  "foo bar $(bas)",
					Start: 0,
					End:   17,
				},
			},
		},
		{
			input: `$("/foo/bar/baz".split('/').slice(-1)[0])`,
			expect: []*Part{
				{
					Raw:   `$("/foo/bar/baz".split('/').slice(-1)[0])`,
					Expr:  `"/foo/bar/baz".split('/').slice(-1)[0]`,
					Start: 0,
					End:   41,
				},
			},
		},
		{
			input: "${\n  var r = [];\n  for (var i = 10; i >= 1; i--) {\n    r.push(i);\n  }\n  return r;\n}\n",
			expect: []*Part{
				{
					Raw:   "${\n  var r = [];\n  for (var i = 10; i >= 1; i--) {\n    r.push(i);\n  }\n  return r;\n}\n",
					Expr:  "var r = [];\n  for (var i = 10; i >= 1; i--) {\n    r.push(i);\n  }\n  return r;",
					Start: 0,
					End:   84,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			t.Logf(`input: "%s"`, test.input)
			parts := Parse(test.input)
			if !reflect.DeepEqual(parts, test.expect) {
				t.Errorf("unexpected matches")
				for _, d := range pretty.Diff(parts, test.expect) {
					t.Log(d)
				}
			}
		})
	}
}

package cwllib

import (
	"github.com/buchanae/cwl"
	"github.com/robertkrimen/otto"
	"regexp"
	"strings"
)

// TODO need parser that tracks open/close of parens
var rx = regexp.MustCompile(`\$\((.*)\)`)

// Part describes a part of a CWL expression string which has been
// parsed by Parse().
type Part struct {
	Raw        string
	Expr       string
	Start, End int
	// true if the expression is a javascript function body (e.g. ${return "foo"})
	IsFuncBody bool
}

// Parse parses a string into a list of parts. If the string does not
// contain a CWL expression, a single part is returned with `Raw` set
// to the original string and `Expr` set to an empty string.
func Parse(expr cwl.Expression) []*Part {
	e := string(expr)
	ev := strings.TrimSpace(e)
	if len(ev) == 0 {
		return nil
	}

	// javascript function expression
	if strings.HasPrefix(ev, "${") && strings.HasSuffix(ev, "}") {
		return []*Part{
			{
				Raw:        e,
				Expr:       strings.TrimSpace(ev[2 : len(ev)-1]),
				Start:      0,
				End:        len(e),
				IsFuncBody: true,
			},
		}
	}

	var parts []*Part

	// parse parameter reference
	last := 0
	matches := rx.FindAllStringSubmatchIndex(e, -1)
	for _, match := range matches {
		start := match[0]
		end := match[1]
		gstart := match[2]
		gend := match[3]

		if start > last {
			parts = append(parts, &Part{
				Raw:   e[last:start],
				Start: last,
				End:   start,
			})
		}

		parts = append(parts, &Part{
			Raw:   string(e[start:end]),
			Expr:  string(e[gstart:gend]),
			Start: start,
			End:   end,
		})
		last = end
	}

	if last < len(e)-1 {
		parts = append(parts, &Part{
			Raw:   string(e[last:]),
			Start: last,
			End:   len(e),
		})
	}

	return parts
}

// IsExpression returns true if the given string contains a CWL expression.
func IsExpression(expr cwl.Expression) bool {
	parts := Parse(expr)
	if len(parts) == 0 {
		return false
	}
	if len(parts) == 1 && parts[0].Expr == "" {
		return false
	}
	return true
}

type ExprData struct {
	Inputs  cwl.Values
	Self    interface{}
	Runtime Runtime
	Libs    []string
}

// Eval evaluates a string which is possibly a CWL expression.
// If the string is not an expression, the string is returned unchanged.
func Eval(e cwl.Expression, data ExprData) (interface{}, error) {
	return EvalParts(Parse(e), data)
}

// EvalParts evaluates a string which has been parsed by Parse().
// If the parts do not represent an expression, the original raw string
// is returned. This is a low-level function, it's better to use EvalString().
func EvalParts(parts []*Part, data ExprData) (interface{}, error) {

	// TODO is there any chance that allocating a VM for every expression is too much?
	//      possibly if used in an API server, this could increase load significantly?
	var vm = otto.New()

	if len(parts) == 0 {
		return nil, nil
	}

	if len(parts) == 1 {
		part := parts[0]

		// No expression, just a normal string.
		if part.Expr == "" {
			return part.Raw, nil
		}

		// Expression or JS function body.
		// Can return any type.
		code := strings.Join(data.Libs, "\n")
		if part.IsFuncBody {
			code = "(function(){" + part.Expr + "})()"
		} else {
			code = "(function(){ return " + part.Expr + "; })()"
		}

		vm.Set("inputs", data.Inputs)
		vm.Set("self", data.Self)
		vm.Set("runtime", map[string]interface{}{
			"outdir":     data.Runtime.Outdir,
			"tmpdir":     data.Runtime.Tmpdir,
			"cores":      data.Runtime.Cores,
			"ram":        data.Runtime.RAM,
			"outdirSize": data.Runtime.OutdirSize,
			"tmpdirSize": data.Runtime.TmpdirSize,
		})

		val, err := vm.Run(code)
		if err != nil {
			return nil, errf("failed to run JS expression: %s", err)
		}

		// otto docs:
		// "Export returns an error, but it will always be nil.
		//  It is present for backwards compatibility."
		ival, _ := val.Export()
		return ival, nil
	}

	// There are multiple parts for expressions of the form "foo $(bar) baz"
	// which is to be treated as string interpolation.

	res := ""
	for _, part := range parts {
		if part.Expr != "" {

			val, err := vm.Run(part.Expr)
			if err != nil {
				return nil, errf("failed to run JS expression: %s", err)
			}

			sval, err := val.ToString()
			if err != nil {
				return nil, errf("failed to convert JS result to a string: %s", err)
			}

			res += sval
		} else {
			res += part.Raw
		}
	}
	return res, nil
}

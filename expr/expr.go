package expr

import (
	"fmt"
	"github.com/kr/pretty"
	"github.com/robertkrimen/otto"
	"regexp"
	"strings"
)

var rx = regexp.MustCompile(`\$\((.*)\)`)

type Part struct {
	Raw        string
	Expr       string
	Start, End int
	// true if the expression is a javascript function body (e.g. ${return "foo"})
	IsFuncBody bool
}

func Parse(e string) []*Part {
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

var vm = otto.New()

func EvalString(s string) (interface{}, error) {
	parts := Parse(s)
	return Eval(parts)
}

// TODO expression results need to go through the loader,
//      so that file types are properly recognized.
func Eval(parts []*Part) (interface{}, error) {
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
		code := part.Expr
		if part.IsFuncBody {
			code = "(function(){" + part.Expr + "})()"
		}

		val, err := vm.Run(code)
		if err != nil {
			return nil, fmt.Errorf("failed to run JS expression: %s", err)
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
				return nil, fmt.Errorf("failed to run JS expression: %s", err)
			}

			sval, err := val.ToString()
			if err != nil {
				return nil, fmt.Errorf("failed to convert JS result to a string: %s", err)
			}

			res += sval
		} else {
			res += part.Raw
		}
	}
	return res, nil
}

func debug(i ...interface{}) {
	pretty.Println(i...)
}

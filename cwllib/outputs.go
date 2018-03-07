package cwllib

import (
	"github.com/buchanae/cwl"
	"github.com/spf13/cast"
	"reflect"
)

/*** CWL output binding code ***/

// Outputs binds cwl.Tool output descriptors to concrete values.
func (job *Job) Outputs() (cwl.Values, error) {
	outdoc, err := job.env.Filesystem().Contents("cwl.output.json")
	if err != nil && err != ErrFileNotFound {
		return nil, err
	}
	if err == nil {
		// TODO type check the output values
		return cwl.LoadValuesBytes([]byte(outdoc))
	}

	values := cwl.Values{}
	for _, out := range job.tool.Outputs {
		v, err := job.bindOutput(out.Type, out.OutputBinding, out.SecondaryFiles, nil)
		if err != nil {
			return nil, errf(`failed to bind value for "%s": %s`, out.ID, err)
		}
		values[out.ID] = v
	}
	return values, nil
}

// bindOutput binds the output value for a single CommandOutput.
func (job *Job) bindOutput(
	types []cwl.OutputType,
	binding *cwl.CommandOutputBinding,
	secondaryFiles []cwl.Expression,
	val interface{},
) (interface{}, error) {
	var err error

	if binding != nil && len(binding.Glob) > 0 {
		// glob patterns may be expressions. evaluate them.
		globs, err := job.evalGlobPatterns(binding.Glob)
		if err != nil {
			return nil, errf("failed to evaluate glob expressions: %s", err)
		}

		files, err := job.matchFiles(globs, binding.LoadContents)
		if err != nil {
			return nil, errf("failed to match files: %s", err)
		}
		val = files
	}

	if binding != nil && binding.OutputEval != "" {
		val, err = job.eval(binding.OutputEval, val)
		if err != nil {
			return nil, errf("failed to evaluate outputEval: %s", err)
		}
	}

	if val == nil {
		for _, t := range types {
			if _, ok := t.(cwl.Null); ok {
				return nil, nil
			}
		}
	}

	if val == nil {
		return nil, errf("missing value")
	}

	debug("value", val)

	// Bind the output value to one of the allowed types.
Loop:
	for _, t := range types {
		switch z := t.(type) {
		case cwl.Boolean:
			v, err := cast.ToBoolE(val)
			if err == nil {
				return v, nil
			}
		case cwl.Int:
			v, err := cast.ToInt32E(val)
			if err == nil {
				return v, nil
			}
		case cwl.Long:
			v, err := cast.ToInt64E(val)
			if err == nil {
				return v, nil
			}
		case cwl.Float:
			v, err := cast.ToFloat32E(val)
			if err == nil {
				return v, nil
			}
		case cwl.Double:
			v, err := cast.ToFloat64E(val)
			if err == nil {
				return v, nil
			}
		case cwl.String:
			v, err := cast.ToStringE(val)
			if err == nil {
				return v, nil
			}
		case cwl.FileType:
			debug("trying to match File")

			switch y := val.(type) {
			case []*cwl.File:
				if len(y) != 1 {
					debug("array is not a single file")
					continue Loop
				}
				f := y[0]
				for _, expr := range secondaryFiles {
					err := job.resolveSecondaryFiles(f, expr)
					if err != nil {
						return nil, errf("failed to resolve secondary files: %s", err)
					}
				}
				return f, nil

			case *cwl.File:
				return y, nil
			case cwl.File:
				return y, nil
			default:
				continue Loop
			}
		case cwl.DirectoryType:
			// TODO
		case cwl.OutputArray:
			debug("trying to match OutputArray")

			typ := reflect.TypeOf(val)
			if typ.Kind() != reflect.Slice {
				debug("value is not an array")
				continue Loop
			}

			var res []interface{}

			arr := reflect.ValueOf(val)
			for i := 0; i < arr.Len(); i++ {
				item := arr.Index(i)
				if !item.CanInterface() {
					return nil, errf("can't get interface of array item")
				}
				r, err := job.bindOutput(z.Items, z.OutputBinding, nil, item.Interface())
				if err != nil {
					return nil, err
				}
				res = append(res, r)
			}
			return res, nil

		case cwl.OutputRecord:
			// TODO
		}
	}

	return nil, errf("no type could be matched")
}

// matchFiles executes the list of glob patterns, returning a list of matched files.
// matchFiles must return a non-nil list on success, even if no files are matched.
func (job *Job) matchFiles(globs []string, loadContents bool) ([]*cwl.File, error) {
	// it's important this slice isn't nil, because the outputEval field
	// expects it to be non-null during expression evaluation.
	files := []*cwl.File{}
	fs := job.env.Filesystem()

	// resolve all the globs into file objects.
	for _, pattern := range globs {
		matches, err := fs.Glob(pattern)
		if err != nil {
			return nil, errf("failed to execute glob: %s", err)
		}

		for _, m := range matches {
			// TODO handle directories
			v := cwl.File{
				Location: m.Location,
				Path:     m.Path,
				Checksum: m.Checksum,
				Size:     m.Size,
			}

			f, err := job.resolveFile(v, loadContents)
			if err != nil {
				return nil, err
			}
			files = append(files, f)
		}
	}
	return files, nil
}

// evalGlobPatterns evaluates a list of potential expressions as defined by the CWL
// OutputBinding.glob field. It returns a list of strings, which are glob expression,
// or an error.
//
// cwl spec:
// "If an expression is provided, the expression must return a string or an array
//  of strings, which will then be evaluated as one or more glob patterns."
func (job *Job) evalGlobPatterns(patterns []cwl.Expression) ([]string, error) {
	var out []string

	for _, pattern := range patterns {
		// TODO what is "self" here?
		val, err := job.eval(pattern, nil)
		if err != nil {
			return nil, err
		}

		switch z := val.(type) {
		case string:
			out = append(out, z)
		case []string:
			out = append(out, z...)
		default:
			return nil, errf(
				"glob expression returned an invalid type. Only string or []string "+
					"are allowed. Got: %s", val)
		}
	}
	return out, nil
}

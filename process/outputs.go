package process

import (
	"github.com/buchanae/cwl"
	"github.com/spf13/cast"
	"reflect"
)

/*** CWL output binding code ***/

// Outputs binds cwl.Tool output descriptors to concrete values.
func (process *Process) Outputs(fs Filesystem) (cwl.Values, error) {
	outdoc, err := fs.Contents("cwl.output.json")
	if err != nil && err != ErrFileNotFound {
		return nil, err
	}
	if err == nil {
		// TODO type check the output values
		return cwl.LoadValuesBytes([]byte(outdoc))
	}

	values := cwl.Values{}
	for _, out := range process.tool.Outputs {
		v, err := process.bindOutput(fs, out.Type, out.OutputBinding, out.SecondaryFiles, nil)
		if err != nil {
			return nil, errf(`failed to bind value for "%s": %s`, out.ID, err)
		}
		values[out.ID] = v
	}
	return values, nil
}

// bindOutput binds the output value for a single CommandOutput.
func (process *Process) bindOutput(
	fs Filesystem,
	types []cwl.OutputType,
	binding *cwl.CommandOutputBinding,
	secondaryFiles []cwl.Expression,
	val interface{},
) (interface{}, error) {
	var err error

	if binding != nil && len(binding.Glob) > 0 {
		// glob patterns may be expressions. evaluate them.
		globs, err := process.evalGlobPatterns(binding.Glob)
		if err != nil {
			return nil, errf("failed to evaluate glob expressions: %s", err)
		}

		files, err := process.matchFiles(fs, globs, binding.LoadContents)
		if err != nil {
			return nil, errf("failed to match files: %s", err)
		}
		val = files
	}

	if binding != nil && binding.OutputEval != "" {
		val, err = process.eval(binding.OutputEval, val)
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

	for _, t := range types {
		switch t.(type) {
		// TODO validate stdout/err can only be at root
		//      validate that stdout/err doesn't occur more than once
		case cwl.Stdout:
			files, err := process.matchFiles(fs, []string{process.stdout}, false)
			if err != nil {
				return nil, errf("failed to match files: %s", err)
			}
			if len(files) == 0 {
				return nil, errf(`failed to match stdout file "%s"`, process.stdout)
			}
			if len(files) > 1 {
				return nil, errf(`matched multiple stdout files "%s"`, process.stdout)
			}
			return files[0], nil

		case cwl.Stderr:
			files, err := process.matchFiles(fs, []string{process.stderr}, false)
			if err != nil {
				return nil, errf("failed to match files: %s", err)
			}
			if len(files) == 0 {
				return nil, errf(`failed to match stderr file "%s"`, process.stderr)
			}
			if len(files) > 1 {
				return nil, errf(`matched multiple stderr files "%s"`, process.stderr)
			}
			return files[0], nil
		}
	}

	if val == nil {
		return nil, errf("missing value")
	}

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
			switch y := val.(type) {
			case []*cwl.File:
				if len(y) != 1 {
					continue Loop
				}
				f := y[0]
				for _, expr := range secondaryFiles {
					err := process.resolveSecondaryFiles(f, expr)
					if err != nil {
						return nil, errf("resolving secondary files: %s", err)
					}
				}
				return f, nil

				// TODO returning both pointer and non-pointer
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
			typ := reflect.TypeOf(val)
			if typ.Kind() != reflect.Slice {
				continue Loop
			}

			var res []interface{}

			arr := reflect.ValueOf(val)
			for i := 0; i < arr.Len(); i++ {
				item := arr.Index(i)
				if !item.CanInterface() {
					return nil, errf("can't get interface of array item")
				}
				r, err := process.bindOutput(fs, z.Items, z.OutputBinding, nil, item.Interface())
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
func (process *Process) matchFiles(fs Filesystem, globs []string, loadContents bool) ([]*cwl.File, error) {
	// it's important this slice isn't nil, because the outputEval field
	// expects it to be non-null during expression evaluation.
	files := []*cwl.File{}

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

			f, err := process.resolveFile(v, loadContents)
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
func (process *Process) evalGlobPatterns(patterns []cwl.Expression) ([]string, error) {
	var out []string

	for _, pattern := range patterns {
		// TODO what is "self" here?
		val, err := process.eval(pattern, nil)
		if err != nil {
			return nil, err
		}

		switch z := val.(type) {
		case string:
			out = append(out, z)
		case []cwl.Value:
			for _, val := range z {
				z, ok := val.(string)
				if !ok {
					return nil, errf(
						"glob expression returned an invalid type. Only string or []string "+
							"are allowed. Got: %#v", z)
				}
				out = append(out, z)
			}
		default:
			return nil, errf(
				"glob expression returned an invalid type. Only string or []string "+
					"are allowed. Got: %#v", z)
		}
	}
	return out, nil
}

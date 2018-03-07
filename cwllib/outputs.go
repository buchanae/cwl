package cwllib

import (
	"github.com/buchanae/cwl"
)

/*** CWL output binding code ***/


// Outputs binds cwl.Tool output descriptors to concrete values.
func (job *Job) Outputs() (cwl.Values, error) {
  // TODO bind the outputs recursively, like bindInput
	values := cwl.Values{}
	for _, out := range job.tool.Outputs {
		v, err := job.bindOutput(out)
		if err != nil {
			return nil, err
		}
		values[out.ID] = v
	}
	return values, nil
}

// bindOutput binds the output value for a single CommandOutput.
func (job *Job) bindOutput(out cwl.CommandOutput) (val interface{}, err error) {
	// glob patterns may be expressions. evaluate them.
	globs, err := job.evalGlobPatterns(out.OutputBinding.Glob)
	if err != nil {
		return nil, errf("failed to evaluate glob expressions for %s: %s", out.ID, err)
	}

	files, err := job.matchFiles(globs, out.OutputBinding.LoadContents)
	if err != nil {
		return nil, errf("failed to match files for %s: %s", out.ID, err)
	}

	if out.OutputBinding.OutputEval != "" {
		val, err := job.eval(out.OutputBinding.OutputEval, files)
		if err != nil {
			return nil, errf("failed to evaluate outputEval for %s: %s", out.ID, err)
		}
		return val, nil
	}
	return files, nil
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

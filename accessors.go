package cwl

import (
	"strings"
)

func (t *Tool) RequiresDocker() (*DockerRequirement, bool) {
	reqs := append([]Requirement{}, t.Requirements...)
	reqs = append(reqs, t.Hints...)
	for _, req := range reqs {
		if r, ok := req.(DockerRequirement); ok {
			return &r, true
		}
	}
	return nil, false
}

func (t *Tool) RequiresShellCommand() bool {
	reqs := append([]Requirement{}, t.Requirements...)
	reqs = append(reqs, t.Hints...)
	for _, req := range reqs {
		if _, ok := req.(ShellCommandRequirement); ok {
			return true
		}
	}
	return false
}

func (t *Tool) RequiresInlineJavascript() ([]string, bool) {
	reqs := append([]Requirement{}, t.Requirements...)
	reqs = append(reqs, t.Hints...)
	for _, req := range reqs {
		if r, ok := req.(InlineJavascriptRequirement); ok {
			return r.ExpressionLib, true
		}
	}
	return nil, false
}

func (t *Tool) RequiresSchemaDef() (*SchemaDefRequirement, bool) {
	reqs := append([]Requirement{}, t.Requirements...)
	reqs = append(reqs, t.Hints...)
	for _, req := range reqs {
		if r, ok := req.(SchemaDefRequirement); ok {
			return &r, true
		}
	}
	return nil, false
}

func (t *Tool) ResolveSchemaDefs() error {
	defs, required := t.RequiresSchemaDef()
	if !required {
		return nil
	}

	byName := map[string]SchemaDef{}
	for _, def := range defs.Types {
		byName[def.Name] = def
	}

	for _, in := range t.Inputs {
		for i, r := range in.Type {
			t, err := resolveSchemaDef(byName, r)
			if err != nil {
				return err
			}
			inputType, ok := t.(InputType)
			if !ok {
				return errf("input type schema reference resolved to non-input type")
			}
			in.Type[i] = inputType
		}
	}
	return nil
}

// resolveSchemaDef recursively resolves references to SchemaDefRequirement types
// by name. Note that schemas defined in the requirement may themselves refer
// to schemas by name. If no type mapping is done, resolveSchemaDef returns the
// original type.
func resolveSchemaDef(byName map[string]SchemaDef, in interface{}) (interface{}, error) {
	switch z := in.(type) {

	case TypeRef:
		name := strings.TrimPrefix(z.Name, "#")
		def, found := byName[name]
		if !found {
			return nil, errf(`no schema def named "%s"`, z.Name)
		}
		// A schema type defined in a SchemaDefRequirement may itself refer
		// to other types by name, so recursively resolve this type.
		//
		// TODO there is a chance of infinite recursion here if types
		//      have a circular reference. Would be nice to catch this with
		//      a good error message.
		t, err := resolveSchemaDef(byName, def.Type.(cwltype))
		if err != nil {
			return nil, err
		}
		return t, nil

	case InputField:
		for i, f := range z.Type {
			t, err := resolveSchemaDef(byName, f)
			if err != nil {
				return nil, err
			}
			inputType, ok := t.(InputType)
			if !ok {
				return nil, errf("input type schema reference resolved to non-input type")
			}
			z.Type[i] = inputType
		}

	case InputRecord:
		for _, f := range z.Fields {
			_, err := resolveSchemaDef(byName, f)
			if err != nil {
				return nil, err
			}
		}

	case InputArray:
		for i, f := range z.Items {
			t, err := resolveSchemaDef(byName, f)
			if err != nil {
				return nil, err
			}
			inputType, ok := t.(InputType)
			if !ok {
				return nil, errf("input type schema reference resolved to non-input type")
			}
			z.Items[i] = inputType
		}
	}
	return in, nil
}

func (clb *CommandLineBinding) GetLoadContents() bool {
	if clb == nil {
		return false
	}
	return clb.LoadContents
}

func (clb *CommandLineBinding) GetPosition() int {
	if clb == nil {
		return 0
	}
	return clb.Position
}

func (clb *CommandLineBinding) GetPrefix() string {
	if clb == nil {
		return ""
	}
	return clb.Prefix
}

func (clb *CommandLineBinding) GetItemSeparator() string {
	if clb == nil {
		return ""
	}
	return clb.ItemSeparator
}

func (clb *CommandLineBinding) GetValueFrom() Expression {
	if clb == nil {
		return Expression("")
	}
	return clb.ValueFrom
}

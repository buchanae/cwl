package cwl

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

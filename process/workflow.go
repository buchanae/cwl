package process

import (
  "fmt"
  "github.com/buchanae/cwl"
)

type scope struct {
  prefix string
  links map[string][]string
}
func (n scope) child(prefix string) scope {
  return scope{prefix: n.prefix + "/" + prefix, links: n.links}
}
func (n scope) link(name string, val string) {
  key := n.key(name)
  n.links[key] = append(n.links[key], val)
}
func (n scope) key(name string) string {
  return n.prefix + "/" + name
}


// DebugWorkflow is a temporary placeholder for workflow processing code.
func DebugWorkflow(wf *cwl.Workflow, vals cwl.Values) {

  root := scope{prefix: "", links: map[string][]string{}}

  inputs := root.child("inputs")
  for k, _ := range vals {
    root.link(k, inputs.key(k))
  }

  exports := linkWorkflow(wf, root)
  outputs := root.child("outputs")

  for _, out := range wf.Outputs {
    outputs.link(out.ID, exports.key(out.ID))
  }

  walk(root.links, []string{"/outputs/count_output"})

}

func walk(links map[string][]string, keys []string) {
  for _, key := range keys {
    fmt.Println(key)
    walk(links, links[key])
  }
}

func linkWorkflow(wf *cwl.Workflow, parent scope) scope {

  internal := parent.child("workflow")
  for _, in := range wf.Inputs {
    internal.link(in.ID, parent.key(in.ID))
  }

  for _, step := range wf.Steps {
    stepScope := internal.child("step/" + step.ID)

    for _, in := range step.In {
      for _, src := range in.Source {
        stepScope.link(in.ID, internal.key(src))
      }
    }

    stepExports := linkDoc(step.Run, stepScope)

    for _, out := range step.Out {
      id := step.ID + "/" + out.ID
      internal.link(id, stepExports.key(out.ID))
    }
  }

  exports := internal.child("exports")
  for _, out := range wf.Outputs {
    for _, src := range out.OutputSource {
      exports.link(out.ID, internal.key(src))
    }
  }

  return exports
}

func linkTool(in []cwl.CommandInput, out []cwl.CommandOutput, parent scope) scope {
  internal := parent.child("tool")
  for _, in := range in {
    internal.link(in.ID, parent.key(in.ID))
    internal.link("toolexec", internal.key(in.ID))
  }
  exports := internal.child("exports")
  for _, out := range out {
    exports.link(out.ID, internal.key("toolexec"))
  }
  return exports
}

func linkDoc(doc cwl.Document, parent scope) scope {
  switch z := doc.(type) {
  case *cwl.Workflow:
    return linkWorkflow(z, parent)
  case *cwl.Tool:
    return linkTool(z.Inputs, z.Outputs, parent)
  case *cwl.ExpressionTool:
    return linkTool(z.Inputs, z.Outputs, parent)
  }
  return scope{}
}

/*
TODO goals

- validate that links are correct, not missing any links, etc
- have (un)marshal-able workflow state
- validate value bindings, mid workflow
- resolve inputs to step in nested workflow, mid workflow
- encode links between steps directly, without intermediate layers


implementation thoughts:
- major element of name translation over many layers
- end result is link between two Process objects and/or
  Step objects.
- possibly want global Start/End steps, or maybe only End;
  End.Done() is true when the workflow is done. need to
  also have link between workflow outputs and last steps
- want to query value of value by name at any layer?
  e.g. query for workflow.step0.count_output mid workflow
*/

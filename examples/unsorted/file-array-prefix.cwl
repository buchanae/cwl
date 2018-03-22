cwlVersion: v1.0
class: CommandLineTool

baseCommand: echo

inputs:
  - id: inputs_separated
    type:
        type: array
        items: File
        inputBinding:
            prefix: "-I"


outputs: []

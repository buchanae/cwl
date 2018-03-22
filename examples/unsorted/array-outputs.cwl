cwlVersion: v1.0
class: CommandLineTool
requirements:
 - class: InlineJavascriptRequirement

baseCommand: echo

arguments:
  - valueFrom: cwl.output.json
    position: 2

inputs:
  touchfiles:
    type: File
    inputBinding:
      position: 1

outputs:
  output:
    type: File[]
    outputBinding:
      glob: "*.txt"
      #outputEval: ${return self}

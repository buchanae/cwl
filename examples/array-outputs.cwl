cwlVersion: v1.0
class: CommandLineTool
requirements:
 - class: InlineJavascriptRequirement
baseCommand: touch
inputs:
  touchfiles:
    type:
      type: array
      items: string
    inputBinding:
      position: 1
outputs:
  output:
    type: File[]
    outputBinding:
      glob: "*.txt"
      #outputEval: ${return self}

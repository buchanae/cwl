cwlVersion: v1.0
class: CommandLineTool
inputs:
  filesA:
    #type: File[]
    type:
      type: array
      items:
        type: array
        items: string
    inputBinding:
      #prefix: -A
      #position: 1

  filesB:
    type:
      type: array
      items: string
      inputBinding:
        prefix: -B=
        separate: false
        #position: -2
    inputBinding:
      prefix: -ZZ
      #position: 20

  filesC:
    type:
      type: array
      items: string
      #inputBinding:
        #prefix: -Z=
        #separate: false
    inputBinding:
      prefix: -C=
      itemSeparator: ","
      separate: false
      position: 4

outputs: []
baseCommand: echo

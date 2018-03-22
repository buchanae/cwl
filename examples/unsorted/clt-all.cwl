cwlVersion: v1.0
class: CommandLineTool
label: Example trivial wrapper for Java 7 compiler
doc: Example doc
hints:
  DockerRequirement:
    dockerPull: java:7-jdk
  DockerRequirement:
    dockerLoad: loadjava:7-jdk
baseCommand:
  - echo
  - foo
arguments: ["-d", $(runtime.outdir)]
stdout: output.txt
stderr: error.txt
inputs:
  scalar: string
  list31: string[]
  list32:
    type: string[]
    inputBinding:
      position: 1
  list4:
    type: string[]
  st1:
    type:
      type: string
      items: string
      bar: baz
  st2:
    type: array?
    items: string
  opt21: string?
  tarfile:
    type: File
    inputBinding:
      position: 1
  extractfile:
    type: string
    inputBinding:
      position: 2
  nullablefile:
    type: ["null", "string"]
    inputBinding:
      position: 2
  list:
    type: string[]
    inputBinding:
      prefix: -A
      position: 3
      itemSeparator: ","
      separate: true
  list2:
    type:
      type: array
      items: string
  optional_file:
    type: File?
  flag:
    type: boolean
  num:
    type: int
outputs:
  scalar1: string
  list3: string[]
  opt2: string?
  output1:
    type: stdout
  error1:
    type: stderr
  example_out:
    type: File
    outputBinding:
      glob: $(inputs.extractfile)
  arrayoutput:
    type:
      type: array
      items: File
  arrayoutput2:
    type: string[]

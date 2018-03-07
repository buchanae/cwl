class: CommandLineTool
cwlVersion: v1.0
requirements:
  - class: ShellCommandRequirement
#hints:
  #DockerRequirement:
    #dockerPull: "debian:wheezy"

inputs: []

outputs:
  - id: output
    type: int

arguments:
   - valueFrom: >
       echo '{"output": "foo" }' > cwl.output.json
     shellQuote: false

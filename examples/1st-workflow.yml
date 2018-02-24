cwlVersion: v1.0
class: Workflow
inputs:
  inp: File
  ex: string
  foo:
    secondaryFiles: .bai
    format: fmt
    doc:
      - doc1
      - doc2
  bar:
    doc: docstring
    format:
      - fm1
      - fm2
    secondaryFiles:
      - .fai
      - .bai

outputs:
  other: File[]
  classout:
    type: File
    outputSource: compile/classfile

steps:
  subwf:
    out: [one]
    in: []
    run:
      class: CommandLineTool
      requirements:
        - class: ShellCommandRequirement
      inputs: []
      arguments:
        - shellQuote: false
          valueFrom: |
            date
            tar cf hello.tar Hello.java
            date
      outputs:
        - id: one
          type: File
          secondaryFiles: .foo
          format: fmt
          outputBinding:
            glob: "*.glob"
          doc:
            - doc1
            - doc2
        - id: arrouttest
          type: File[]
          doc: docstring
          outputBinding:
            glob:
              - "*.glob1"
              - "*.glob2"
          format:
            - fm1
            - fm2
          secondaryFiles:
            - .fai
            - .bai

  auntar:
    run: tar-param.cwl
    in: 
      tarfile: inp
      other:
        source:
          - ex
      extractfile: ex
    out:
      - example_out
    scatter: tarfile
  untar:
    run: tar-param.cwl
    in:
      tarfile: inp
      extractfile: ex
    out:
      - id: example_out

  compile:
    run: arguments.cwl
    in:
      src: untar/example_out
    out: [classfile]

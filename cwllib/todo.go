package cwllib

/*
TODO
- more complete JS expression context (self, inputs, runtime, etc)
- output document binding
- cwl.output.json
- secondary files
- load expression result values into File/Directory types where appropriate
- file staging and working directory
- solid expression parser (regexp misses edge cases and escaping)
- relative path context (current working directory) for filesystems
- absolute paths for files, especially in outputs
- resolve document references
- filesystem multiplexing based on location
- success/failure codes and relationship to CLI cmd
- Any type
- document validation before processing
- better line/col info from document loading errors
- carefully check document json/yaml marshaling
- input/output record type handling
- executor backends
- directory type
- good framework for e2e tests with lots of coverage
- $include and $import
- test unrecognized fields are ignored (possibly with warning)
- optional checksum calculation for filesystems
- resource requests
- environment variables
- initial work dir
- docker
- missing requirement/hint types. see requirements.go
- time limit on JS evaluation

workflow execution:
- basics
- caching

server + API:

cmd/plan:
  "cwl plan" command which describes (in JSON?) what will happen when executing
  a workflow or tool, including file transfers etc.
*/

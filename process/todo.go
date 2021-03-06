package process

/*
TODO
- load expression result values into File/Directory types where appropriate
- absolute paths for files, especially in outputs
- good framework for e2e tests with lots of coverage
- really good debug logging, with the goal of clearly explaining to a **user**
  what is going on when a job fails at any step, especially input/output binding.
- success/failure codes and relationship to CLI cmd
- solid expression parser (regexp misses edge cases and escaping)
- type check cwl.output.json
- filesystem multiplexing based on location

- document validation before processing
- better line/col/context info from document loading errors
- carefully check document json/yaml marshaling
- input/output record type handling
- executor backends
- directory type
- $include and $import
- test unrecognized fields are ignored (possibly with warning)
- optional checksum calculation for filesystems
- resource requests
- initial work dir
- time limit on JS evaluation

workflow execution:
- basics
- caching

server + API:

cmd/plan:
  "cwl plan" command which describes (in JSON?) what will happen when executing
  a workflow or tool, including file transfers etc.
*/

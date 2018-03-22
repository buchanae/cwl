cwltest () {
  cat $1-*/meta.yaml
  echo
  cwl run $1-*/tool.cwl $1-*/job.cwl
}

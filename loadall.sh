for f in $(ls ~/src/github.com/kids-first/kf-alignment-workflow/tools/*.cwl); do
  if cwl dump $f > /dev/null; then
    echo $f
  else
    echo $f
  fi
done

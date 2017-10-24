package cwl

import (
  "fmt"
  "testing"
  "io/ioutil"
	"github.com/kr/pretty"
)

func TestLoad1stTool(t *testing.T) {
  b, _ := ioutil.ReadFile("examples/1st-tool.yml")
  fmt.Println(string(b))
  tool, err := ParseTool(string(b))
  pretty.Println(tool)
  pretty.Println(err)
}

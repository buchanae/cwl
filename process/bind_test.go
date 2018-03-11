package process

/* TODO
import (
	"reflect"
	"testing"
)

func TestFormatArgs(t *testing.T) {
	clbA := &CommandLineBinding{
		Prefix: "-A",
	}

	clbB := &CommandLineBinding{
		Prefix: "-B=",
	}
	clbB.SetSeparate(false)

	fa := formatArgs(clbA, "foo")
	xa := []string{"-A", "foo"}

	fb := formatArgs(clbB, "foo")
	xb := []string{"-B=foo"}

	if !reflect.DeepEqual(fa, xa) {
		t.Error("incorrect formatting, should be separated by default")
	}

	if !reflect.DeepEqual(fb, xb) {
		t.Error("incorrect formatting, separate is false")
	}
}
*/

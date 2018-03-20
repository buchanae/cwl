package local

import (
	"fmt"
)

// errf makes fmt.Errorf shorter
func errf(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}

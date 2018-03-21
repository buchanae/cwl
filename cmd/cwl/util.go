package main

import (
	"fmt"
	"github.com/kr/pretty"
)

// errf makes fmt.Errorf shorter
func errf(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}

func debug(i ...interface{}) {
	pretty.Println(i...)
}

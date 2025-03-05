package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/aeremic/cgo/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome %s to C Go programming langauge!\n", user.Name)
	repl.Start(os.Stdin, os.Stdout)
}

package main

import (
	"fmt"
	"os"

	"github.com/lincolnauster/painted/pkg"
)

func exit(e error) {
	fmt.Fprintln(os.Stderr, e.Error())
	os.Exit(1)
}

func main() {
	args := painted.DefaultArgs()

	if err := args.Fill(os.Args[1:]); err != nil {
		exit(err)
	}

	action, err := args.Make()
	if err != nil {
		exit(err)
	}

	err = action.Exec()
	if err != nil {
		exit(err)
	}
}

package main

import (
	"log"

	"github.com/lincolnauster/painted/pkg/painted"
)

func main() {
	args, err := painted.FromArgs()

	if err != nil {
		log.Fatal(err)
	}

	action, err := args.Make()
	if err != nil {
		log.Fatal(err)
	}

	err = action.Exec()
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"os"

	"github.com/lincolnauster/painted/pkg/painted"
)

func main() {
	args := painted.DefaultArgs()

	if err := args.Fill(os.Args[1:]); err != nil {
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

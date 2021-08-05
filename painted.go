package main

import "os"

func main() {
	args := Args{
		false,
		"/dev/stdin",
		"/dev/stdout",
	}

	if err := args.Fill(os.Args[1:]); err != nil {
		panic(err)
	}

	action, err := args.Make()
	if err != nil {
		panic(err)
	}

	err = action.Exec()
	if err != nil {
		panic(err)
	}
}

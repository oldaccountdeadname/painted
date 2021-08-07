package main

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"net"
	"os"

	"gitlab.com/lincolnauster/painted/dbus"
)

var HelpMessage = Out{
	`painted

  usage: painted [-h|[--]help] | [-i|[--]input path] [-o|[--]output path]

If a path ends in .sock, it is interpreted as a UNIX socket.
`,
}

var (
	AvailableArgs = map[byte]string{
		'h': "help",
		'i': "input",
		'o': "output",
	}

	ArgRequiresVal = map[string]bool{
		"help":   false,
		"input":  true,
		"output": true,
	}
)

type Exec interface {
	Exec() error
}

type Args struct {
	Help   bool
	Input  string
	Output string
}

type Out struct {
	msg string
}

// Apply a CLI arg given a key and an associated value (the latter of which may
// be nil). If arguments are invalid, they are ignored.
func (a *Args) Apply(k *string, v *string) {
	if *k == "help" {
		a.Help = true
	} else if v == nil {
		return
	}

	if *k == "input" {
		a.Input = *v
	} else if *k == "output" {
		a.Output = *v
	}
}

// Initialize `self` with a list of arguments. Errors returned are fatal.
func (a *Args) Fill(args_s []string) error {
	for i := 0; i < len(args_s); i++ {
		var val string
		key, m := argToOpt(args_s[i])

		if m == true {
			return errors.New(fmt.Sprintf(
				"Arg %s is malformed. Must begin -- or -",
				key,
			))
		}

		if ArgRequiresVal[key] {
			i += 1

			if i >= len(args_s) {
				return errors.New(
					fmt.Sprintf(
						"Arg %s requires value.",
						key,
					),
				)
			}

			val = args_s[i]
		}

		a.Apply(&key, &val)
	}

	return nil
}

func (a *Args) Make() (Exec, error) {
	if a.Help {
		return HelpMessage, nil
	} else {
		reader, r_err := asReader(a.Input)
		writer, w_err := asWriter(a.Output)

		e_msg := ""
		if r_err != nil {
			e_msg += fmt.Sprintf(
				"Error opening file %s for reading: %s\n",
				a.Input, r_err.Error(),
			)
		}
		if w_err != nil {
			e_msg += fmt.Sprintf(
				"Error opening file %s for writing: %s\n",
				a.Output, w_err.Error(),
			)
		}

		if e_msg != "" {
			return nil, errors.New(e_msg)
		} else {
			return Model{
				a.Input,
				reader,
				writer,
				dbus.SessionConn{nil},
			}, nil
		}
	}
}

func (o Out) Exec() error {
	fmt.Printf("%s", o.msg)
	return nil
}
func asReader(p string) (io.Reader, error) {
	if strings.HasSuffix(p, ".sock") {
		a, b := net.Dial("unix", p)
		return a, b
	} else {
		f, e := os.OpenFile(p, os.O_RDONLY, 0664)
		if e == nil {
			// seek to the end and don't reread old commands
			f.Seek(0, 2)
		}
		return f, e
	}
		
}

func asWriter(p string) (io.Writer, error) {
	return os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 664)
}

// reduce a CLI arg from --string or -s to string (or the corresponding version
// of the short representation. If neither pattern matches, the second return
// value is `true` to indicate a malformed argument.
func argToOpt(s string) (string, bool) {
	if s[:(2%len(s))] == "--" {
		return s[2:], false
	} else if s[0] == '-' {
		return AvailableArgs[s[1]], false
	} else {
		return s, true
	}

}

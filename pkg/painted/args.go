package painted

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/gammazero/deque"

	"github.com/lincolnauster/painted/pkg/dbus"
)

var HelpMessage = Out{
	`painted

  usage: painted --help] | --input <path> --output <path> --config <path>

If input or output paths end in .sock, they're interpreted as a UNIX socket.
`,
}

var (
	AvailableArgs = map[byte]string{
		'h': "help",
		'i': "input",
		'o': "output",
		'c': "config",
	}

	ArgRequiresVal = map[string]bool{
		"help":   false,
		"input":  true,
		"output": true,
		"config": true,
	}
)

type Exec interface {
	Exec() error
}

type Args struct {
	Help   bool
	Input  string
	Output string
	Config string
}

type Out struct {
	msg string
}

func FromArgs() (*Args, error) {
	args := defaultArgs()

	if err := args.fill(os.Args[1:]); err != nil {
		return nil, err
	}

	return &args, nil
}

func defaultArgs() Args {
	conf_location := os.Getenv("HOME") + "/.config/painted/conf.toml"
	return Args{
		false,
		"/dev/stdin",
		"/dev/stdout",
		conf_location,
	}
}

// Initialize `self` with a list of arguments. Errors returned are fatal.
func (a *Args) fill(args_s []string) error {
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

		a.apply(&key, &val)
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
			conf, _ := MakeConfigFromFile(a.Config)
			return Model{
				conf,
				Io{
					Reader{reader, a.Input},
					Writer{writer, a.Output},
				},
				dbus.SessionConn{},
				NotifQueue{*deque.New(), 0},
			}, nil
		}
	}
}

func (o Out) Exec() error {
	fmt.Printf("%s", o.msg)
	return nil
}

// Apply a CLI arg given a key and an associated value (the latter of which may
// be nil). If arguments are invalid, they are ignored.
func (a *Args) apply(k *string, v *string) {
	if *k == "help" {
		a.Help = true
	}

	if *k == "input" && v != nil {
		a.Input = *v
	} else if *k == "output" && v != nil {
		a.Output = *v
	}
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
	return os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
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

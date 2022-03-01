package painted

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/gammazero/deque"
)

type Exec interface {
	Exec() error
}

type Args struct {
	Input  string
	Output string
	Config string
}

type Out struct {
	msg string
}

func FromArgs() (*Args, error) {
	var args Args

	flag.StringVar(
		&args.Input,
		"input",
		"/dev/stdin",
		"The input file to read from.",
	)

	flag.StringVar(
		&args.Output,
		"output",
		"/dev/stdout",
		"The output file to read from.",
	)

	flag.StringVar(
		&args.Config,
		"config",
		os.Getenv("HOME")+"/.config/painted/conf.toml",
		"The config file.",
	)

	flag.Parse()

	return &args, nil
}

func (a *Args) Make() (Exec, error) {
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
		conf, err := MakeConfigFromFile(a.Config)
		if err != nil {
			return nil, err
		}

		return Model{
			conf,
			Io{
				Reader{reader, a.Input},
				Writer{writer, a.Output},
			},
			NotifQueue{*deque.New(), 0},
		}, nil
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
		f, e := os.OpenFile(p, os.O_RDONLY|os.O_CREATE, 0664)
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

//+build freebsd || openbsd

package painted

import (
	"bufio"
	"fmt"
	"io"
	"syscall"
)

type Reader struct {
	File io.Reader
	Path string
}

type Writer struct {
	File io.Writer
	Path string
}

type Io struct {
	Reader Reader
	Writer Writer
}

/* reading helpers */

// Return a closure that continues to read the next line. On EOF, the returned
// function blocks until the file has been modified, and then tries again.
func (i *Io) Lines() func() (string, error) {
	f := bufio.NewReader(i.Reader.File)

	return func() (string, error) {
		for {
			out, err := f.ReadString('\n')
			if err == io.EOF {
				blockUntilModify(i.Reader.Path)
			} else if err != nil {
				return "", err
			} else {
				return out, nil
			}
		}
	}
}

/* writing helpers */

func (i *Io) Write(s string) {
	i.Writer.File.Write([]byte(s))
}

func (i *Io) Writef(f string, v ...interface{}) {
	i.Write(fmt.Sprintf(f, v...))
}

/* misc helpers */

// Use kqueue to watch a given file path and block until an event occurs.
//
// Errors are ignored.
func blockUntilModify(f string) {
	kq , _ := syscall.Kqueue()
	fd, _ := syscall.Open(f, syscall.O_RDONLY, 0)
	// seek to the end to only register new writes
	// FIXME: this isn't atomic.
	syscall.Seek(fd, 0, 2)

	ke := syscall.Kevent_t{
		Ident: uint64(fd),
		Filter: syscall.EVFILT_READ,
		Flags: syscall.EV_ADD | syscall.EV_ONESHOT,
		Fflags: 0,
		Data: 0,
		Udata: nil,
	}

	kchange_list := []syscall.Kevent_t{ke}
	kevent_list  := make([]syscall.Kevent_t, 1);

	syscall.Kevent(kq, kchange_list, kevent_list, nil)
}

package painted

import (
	"bufio"
	"fmt"
	"io"
	"unsafe"

	"golang.org/x/sys/unix"
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

// Use inotify to watch a given file path and `read` (block until an event
// occurs). See inotify(2). This is a linux-specific syscall.
//
// Errors are ignored.
func blockUntilModify(f string) {
	nf, err := unix.InotifyInit()

	if err != nil {
		return
	}

	unix.InotifyAddWatch(nf, f, unix.IN_MODIFY)

	ev := make([]byte, unsafe.Sizeof(unix.InotifyEvent{}))
	unix.Read(nf, ev)
}

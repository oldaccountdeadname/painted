//go:build linux
// +build linux

package painted

import (
	"unsafe"

	"golang.org/x/sys/unix"
)

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

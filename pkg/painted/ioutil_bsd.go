//go:build freebsd || openbsd
// +build freebsd openbsd

package painted

import "syscall"

// Use kqueue to watch a given file path and block until an event occurs.
//
// Errors are ignored.
func blockUntilModify(f string) {
	kq, _ := syscall.Kqueue()
	fd, _ := syscall.Open(f, syscall.O_RDONLY, 0)
	// seek to the end to only register new writes
	// FIXME: this isn't atomic.
	syscall.Seek(fd, 0, 2)

	ke := syscall.Kevent_t{
		Ident:  uint64(fd),
		Filter: syscall.EVFILT_READ,
		Flags:  syscall.EV_ADD | syscall.EV_ONESHOT,
		Fflags: 0,
		Data:   0,
		Udata:  nil,
	}

	kchange_list := []syscall.Kevent_t{ke}
	kevent_list := make([]syscall.Kevent_t, 1)

	syscall.Kevent(kq, kchange_list, kevent_list, nil)
}

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"sync/atomic"
	"unsafe"

	"gitlab.com/lincolnauster/painted/dbus"
	"golang.org/x/sys/unix"
)

// The model links together dbus and IO interaction into one entry point.
type Model struct {
	InputPath  string
	InputFile  io.Reader
	OutputFile io.Writer
	bus        dbus.SessionConn
}

// This structure implements dbus' org.freedesktop.Notifications interface and
// encapsulates state. It's useful as an object to be exported onto the session
// bus at /org/freedesktop/Notifications.
type Server struct {
	NextId uint32
	Model  *Model
}

// This is an in-memory representation of the notification for manipulation onto
// IO. It is *not* a direct mapping of the notification spec[0] and contains
// only the information that is used by painted.
//
// [0]: https://developer-old.gnome.org/notification-spec/
type Notification struct {
	OriginApp string
	Summary   string
	Id        uint32
	ReplaceId uint32
	/* if Id == ReplaceId, then the notif doesn't replace anything. */
}

func (m *Model) takeName() error {
	reply := m.bus.TakeName(
		"org.freedesktop.Notifications",
	)
	if reply != true {
		return errors.New(
			`Can't take org.freedesktop.Notifications. Is another notif daemon running?`,
		)
	}

	return nil
}

func (m *Model) RegisterIface(serv *Server) error {
	if err := m.bus.Export(
		serv,
		"/org/freedesktop/Notifications",
		"org.freedesktop.Notifications",
	); err != nil {
		return err
	} else {
		return nil
	}
}

// Continuously read lines from a file. This does *not* respect EOF, and behaves
// similarly to `tail -f`. (TODO)
func (m *Model) CmdLoop() {
	f := bufio.NewReader(m.InputFile)
	for {
		cmd, err := f.ReadString('\n')
		if err == io.EOF {
			blockUntilModify(m.InputPath)
		} else if err != nil {
			return
		} else {
			// Otherwise, cmd is perfectly valid.

			// strip the newline
			cmd = cmd[:len(cmd)-1]

			// TODO *way* better matching logic
			// (thinking a trie for prefix-matching)
			switch cmd {
			case "exit":
				return
			case "clear":
				m.OutputFile.Write([]byte{'\n'})
			default:
				m.OutputFile.Write([]byte(
					fmt.Sprintf("%s not understood.\n", cmd),
				))
			}
		}
	}
}

func (m *Model) Notify(n Notification) {
	// TODO pretty formattting
	m.OutputFile.Write([]byte(fmt.Sprintf("%+v\n", n)))
}

// Connect to the bus, register the interface, launch the notif loop and the
// input loop (concurrently).
func (m Model) Exec() error {
	defer m.bus.Close()
	
	if err := m.takeName(); err != nil {
		return err
	}

	var serv Server
	serv.Model = &m
	serv.NextId = 1

	if err := m.RegisterIface(&serv); err != nil {
		return err
	}

	m.CmdLoop()

	return nil
}

func (s *Server) GetServerInformation() (
	string, string, string, string, *dbus.Error,
) {
	return "painted", "none", "v0.1.0", "v1.2", nil
}

func (s *Server) GetCapabilities() ([]string, *dbus.Error) {
	fmt.Println("GetCapabilities called.")
	return []string{}, nil
}

func (s *Server) Notify(
	app_name string,
	replaces_id uint32,
	app_icon string,
	summary string,
	body string,
	actions []interface{},
	hints map[interface{}]interface{},
	expire_timeout int32,
) (uint32, *dbus.Error) {
	notif := Notification{
		OriginApp: app_name,
		Summary:   summary,
		Id:        s.NextId,
		ReplaceId: replaces_id,
	}

	// If the sender doesn't want to replace a notification, they send 0 as
	// the replace ID. If that's the case, we assign the replace ID the
	// newly allocated id to indicate that this notification is not a
	// replacement for anything, and then we have to update nextID
	// accordingly. Otherwise, we can recycle the ID.
	if notif.ReplaceId == 0 {
		notif.ReplaceId = notif.Id
		atomic.AddUint32(&s.NextId, 1)
	} else {
		notif.Id = notif.ReplaceId
	}

	s.Model.Notify(notif)

	return notif.Id, nil
}

// Use inotify to watch a given file path and `read` (block until an event
// occurs). See inotify(2). This is a linux-specific syscall.
//
// The wait is level-triggered so as not to block indefinitely if called in a
// loop.
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

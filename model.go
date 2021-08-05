package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	"github.com/godbus/dbus/v5"
)

type Model struct {
	inputF     io.Reader
	OutputFile io.Writer
	Bus        *dbus.Conn
}

// This structure implements dbus' org.freedesktop.Notifications interface and
// encapsulates state. It's useful as an object to be exported onto the session
// bus at /org/freedesktop/Notifications.
type Server struct {
	nextId uint32
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
	ReplaceId *uint32
}

// Connect to the bus. If an error occurs, m.Bus is set to nil.
func (m *Model) connect() error {
	bus, err := dbus.ConnectSessionBus()
	m.Bus = bus
	return err
}

func (m *Model) takeName() error {
	reply, err := m.Bus.RequestName(
		"org.freedesktop.Notifications",
		dbus.NameFlagReplaceExisting,
	)

	if err != nil {
		return err
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		return errors.New(
			`Can't take org.freedesktop.Notifications. Is another notif daemon running?`,
		)
	}

	return nil
}

func (m *Model) releaseName() {
	m.Bus.ReleaseName("org.freedesktop.Notifications")
}

func (m *Model) RegisterIface(serv *Server) error {
	m.Bus.BusObject().AddMatchSignal(
		"org.freedesktop.Notifications",
		"GetServerInformation",
	)

	if err := m.Bus.Export(
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
	f := bufio.NewReader(m.inputF)
	for {
		cmd, err := f.ReadString('\n')
		if err == io.EOF {
			// TODO block until modification, then continue.

			// right now, it just spins, which is sorta less than
			// ideal.
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
	if err := m.connect(); err != nil {
		return err
	} else {
		defer m.Bus.Close()
	}

	if err := m.takeName(); err != nil {
		return err
	} else {
		defer m.releaseName()
	}

	var serv Server
	serv.Model = &m

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
		Id:        s.nextId,
		ReplaceId: &replaces_id,
	}

	// only replace if the sender indicated it
	if *notif.ReplaceId == 0 {
		notif.ReplaceId = nil
	}

	s.nextId += 1

	s.Model.Notify(notif)

	return notif.Id, nil
}

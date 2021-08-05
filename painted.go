package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/godbus/dbus/v5"
)

var (
	bufSize = 8
)

type Model struct {
	inputF     os.File
	outputF    os.File
	Bus        *dbus.Conn
	WriteQueue *chan string
}

// This structure implements dbus' org.freedesktop.Notifications interface and
// encapsulates state. It's useful as an object to be exported onto the session
// bus at /org/freedesktop/Notifications.
type Server struct {
	nextId     uint32
	WriteQueue *chan string
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

	if m.WriteQueue == nil {
		ch := make(chan string, 16)
		m.WriteQueue = &ch
	}

	var serv Server
	serv.WriteQueue = m.WriteQueue

	if err := m.RegisterIface(&serv); err != nil {
		return err
	}

	for {
		m.outputF.WriteString(<- *m.WriteQueue)
	}

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

	*s.WriteQueue <- fmt.Sprintf("%+v\n", notif) // TODO pretty formatting

	return notif.Id, nil
}

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

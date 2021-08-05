package main

import (
	"os"
	"fmt"
	"errors"

	"github.com/godbus/dbus/v5"
)

var (
	bufSize = 8
)

type Model struct {
	Input  os.File
	Output os.File
	Bus    *dbus.Conn
}

// We need a 'dummy' struct of sorts to export out into dbus. It provides the
// interface that org.freedesktop.Notifications is expecting.
type Server struct {}

// Connect. If an error occurs, m.Bus is set to nil.
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
			"unable to become primary owner of org.freedesktop.Notifications.",
		)
	}

	return nil
}

func (m *Model) releaseName() {
	m.Bus.ReleaseName("org.freedesktop.Notifications")
}


func (m *Model) RegisterIface() error {
	m.Bus.BusObject().AddMatchSignal(
		"org.freedesktop.Notifications",
		"GetServerInformation",
	)

	if err := m.Bus.Export(
		new(Server),
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

	if err := m.RegisterIface(); err != nil {
		return err
	}

	for {

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
) *dbus.Error {
	fmt.Printf("%s\n", summary)
	return nil
}

func main() {
	args := Args{
		false,
		"/dev/stdin",
		"/dev/stdout",
	}

	args.Fill(os.Args[1:])

	action, err := args.Make()
	if err != nil {
		panic(err)
	}
	
	err = action.Exec()
	if err != nil {
		panic(err)
	}
}

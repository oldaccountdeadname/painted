package dbus

import "github.com/godbus/dbus/v5"

// This wraps the go dbus lib and provides an idiomatic interface. It abstracts
// over things like connection state by providing a lazy API.
type SessionConn struct {
 	Conn *dbus.Conn
}

type Error = dbus.Error

// Take a name. If the name is not immediately provided or another error occurs,
// `false` is returned. This will *not* queue, and will replace if possible. It
// will not allow replacement itself.
//
// Names are automatically returned to dbus when the connection is closed.
func (s *SessionConn) TakeName(name string) bool {
	s.lazyConnect()

	resp, err := s.Conn.RequestName(
		name,
		dbus.NameFlagReplaceExisting | dbus.NameFlagDoNotQueue,
	)

	return err == nil && resp == dbus.RequestNameReplyPrimaryOwner
}

func (s *SessionConn) Export(obj interface{}, path string, iface string) error {
	return s.Conn.Export(obj, dbus.ObjectPath(path), iface)
}

func (s *SessionConn) Close() {
	s.Conn.Close()
	s.Conn = nil
}

func (s *SessionConn) lazyConnect() error {
	var conn *dbus.Conn
	var err  error
	if s.Conn == nil {
		conn, err = dbus.ConnectSessionBus()
	} else {
		conn = s.Conn
		err  = nil
	}

	s.Conn = conn
	return err
}


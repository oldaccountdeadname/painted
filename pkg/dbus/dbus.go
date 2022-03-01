package dbus

import "github.com/godbus/dbus/v5"

// This wraps the go dbus lib and provides an idiomatic interface. It abstracts
// over things like connection state by providing a lazy API.
type SessionConn struct {
	conn *dbus.Conn
}

type Error = dbus.Error

// Take a name. If the name is not immediately provided or another error occurs,
// `false` is returned. This will *not* queue, and will replace if possible. It
// will not allow replacement itself.
//
// Names are automatically returned to dbus when the connection is closed.
func (s *SessionConn) TakeName(name string) bool {
	s.lazyConnect()

	resp, err := s.conn.RequestName(
		name,
		dbus.NameFlagReplaceExisting|dbus.NameFlagDoNotQueue,
	)

	return err == nil && resp == dbus.RequestNameReplyPrimaryOwner
}

func (s *SessionConn) Export(obj interface{}, path string, iface string) error {
	return s.conn.Export(obj, dbus.ObjectPath(path), iface)
}

func (s *SessionConn) Emit(
	path dbus.ObjectPath,
	name string,
	values ...interface{},
) error {
	return s.conn.Emit(path, name, values...)
}

func (s *SessionConn) Close() {
	s.conn.Close()
	s.conn = nil
}

func (s *SessionConn) lazyConnect() error {
	var conn *dbus.Conn
	var err error
	if s.conn == nil {
		conn, err = dbus.ConnectSessionBus()
	} else {
		conn = s.conn
		err = nil
	}

	s.conn = conn
	return err
}

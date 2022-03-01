package dbus

import "github.com/godbus/dbus/v5"

// This wraps the go dbus lib and provides an idiomatic interface. It abstracts
// over things like connection state by providing a lazy API.

var conn *dbus.Conn

type Error = dbus.Error

// Take a name. If the name is not immediately provided or another error occurs,
// `false` is returned. This will *not* queue, and will replace if possible. It
// will not allow replacement itself.
//
// Names are automatically returned to dbus when the connection is closed.
func TakeName(name string) bool {
	lazyConnect()

	resp, err := conn.RequestName(
		name,
		dbus.NameFlagReplaceExisting|dbus.NameFlagDoNotQueue,
	)

	return err == nil && resp == dbus.RequestNameReplyPrimaryOwner
}

func Export(obj interface{}, path string, iface string) error {
	return conn.Export(obj, dbus.ObjectPath(path), iface)
}

func Emit(
	path dbus.ObjectPath,
	name string,
	values ...interface{},
) error {
	return conn.Emit(path, name, values...)
}

func Close() {
	conn.Close()
	conn = nil
}

func lazyConnect() error {
	var err error
	if conn == nil {
		conn, err = dbus.ConnectSessionBus()
		return err
	} else {
		return nil
	}
}

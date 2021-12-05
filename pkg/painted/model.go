package painted

import (
	"errors"
	"sync/atomic"

	"github.com/lincolnauster/painted/pkg/dbus"
	"github.com/lincolnauster/painted/pkg/trie"
)

// The model links together dbus and IO interaction into one entry point.
type Model struct {
	conf  Config
	io    Io
	bus   dbus.SessionConn
	queue NotifQueue
}

// This structure implements dbus' org.freedesktop.Notifications interface and
// encapsulates state. It's useful as an object to be exported onto the session
// bus at /org/freedesktop/Notifications.
type listener struct {
	nextId  uint32
	Recieve func(Notification)
}

func (m *Model) takeName() error {
	if !m.bus.TakeName("org.freedesktop.Notifications") {
		return errors.New(
			`Can't take org.freedesktop.Notifications. Is another notif daemon running?`,
		)
	} else {
		return nil
	}
}

func (m *Model) registerIface(listener *listener) error {
	return m.bus.Export(
		listener,
		"/org/freedesktop/Notifications",
		"org.freedesktop.Notifications",
	)
}

// Assuming a properly set up Io member, read and process all input commands as
// they arise (irrespective of EOF). This is blocking until the exit command is
// sent.
//
// TODO: extract this into a goroutine managed by a channel which is a member of
// this struct. Then performCmd may *actually perform the command* by sending
// any IO changes (exit) across that channel.
func (m *Model) CmdLoop() {
	var cmd_trie trie.Trie
	cmd_trie.Insert([]rune("exit"))
	cmd_trie.Insert([]rune("clear"))
	cmd_trie.Insert([]rune("next"))
	cmd_trie.Insert([]rune("previous"))
	cmd_trie.Insert([]rune("expand"))
	cmd_trie.Insert([]rune("help"))

	next_line := m.io.Lines()

	for {
		cmd, err := next_line()

		if err != nil {
			panic(err)
		}

		cmd = cmd[:len(cmd)-1]
		cmd_r := []rune(cmd)
		match := cmd_trie.SearchWithDefault(cmd_r, cmd_r)

		if m.performCmd(string(match)) {
			break
		}
	}
}

func (m *Model) Notify(n Notification) {
	m.queue.Push(&n)
	m.queue.CallOnCurrent(func(n *Notification) {
		m.io.Write(m.conf.Formatter(n))
	})
}

// Connect to the bus, register the interface, launch the notif loop and the
// input loop (concurrently).
func (m Model) Exec() error {
	defer m.bus.Close()

	if err := m.takeName(); err != nil {
		return err
	}

	var listener listener
	listener.Recieve = m.Notify

	if err := m.registerIface(&listener); err != nil {
		return err
	}

	m.CmdLoop()

	return nil
}

// Perform a command, and return true if an exit was requested.
func (m *Model) performCmd(cmd string) bool {
	switch cmd {
	case "exit":
		return true
	case "clear":
		m.io.Write("\n")
	case "next":
		m.queue.Next()
		m.queue.CallOnCurrent(func(n *Notification) {
			m.io.Write(m.conf.Formatter(n))
		})
	case "previous":
		m.queue.Prev()
		m.queue.CallOnCurrent(func(n *Notification) {
			m.io.Write(m.conf.Formatter(n))
		})
	case "expand":
		m.queue.CallOnCurrent(func(n *Notification) {
			m.io.Writef("%s\n", n.Body)
		})
	case "help":
		m.io.Write(
			"command should be: exit | clear | next | previous | help\n",
		)
	default:
		m.io.Writef("%s not matched with any valid commands: see `help`.\n", cmd)
	}

	return false
}

func (l *listener) GetServerInformation() (
	string, string, string, string, *dbus.Error,
) {
	return "painted", "none", "v0.1.1", "v1.2", nil
}

func (l *listener) GetCapabilities() ([]string, *dbus.Error) {
	return []string{"persistence", "body"}, nil
}

func (l *listener) Notify(
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
		Body:      body,
		Id:        replaces_id,
	}

	if notif.Id == 0 { // test if we need to give this an ID
		notif.Id = atomic.AddUint32(&l.nextId, 1)
	}

	l.Recieve(notif)

	return notif.Id, nil
}

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
	if !dbus.TakeName("org.freedesktop.Notifications") {
		return errors.New(
			`Can't take org.freedesktop.Notifications. Is another notif daemon running?`,
		)
	} else {
		return nil
	}
}

func (m *Model) registerIface(listener *listener) error {
	return dbus.Export(
		listener,
		"/org/freedesktop/Notifications",
		"org.freedesktop.Notifications",
	)
}

// Assuming a properly set up Io member, read and process all input commands as
// they arise (irrespective of EOF). This is blocking until the exit command is
// sent.
func (m *Model) CmdLoop() {
	var cmd_trie trie.Trie
	cmd_trie.Insert([]rune("exit"))
	cmd_trie.Insert([]rune("clear"))
	cmd_trie.Insert([]rune("remove"))
	cmd_trie.Insert([]rune("next"))
	cmd_trie.Insert([]rune("previous"))
	cmd_trie.Insert([]rune("expand"))
	cmd_trie.Insert([]rune("summarize"))
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
	m.summarizeNotif()
}

// Connect to the bus, register the interface, launch the notif loop and the
// input loop (concurrently).
func (m Model) Exec() error {
	defer dbus.Close()

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
		m.queue.Get().Dismiss()
		m.io.Write("\n")
	case "remove":
		m.performCmd("clear")
		m.queue.Remove().Dismiss()
	case "next":
		m.queue.Next()
		m.summarizeNotif()
	case "previous":
		m.queue.Prev()
		m.summarizeNotif()
	case "expand":
		m.expandNotif()
	case "summarize":
		m.summarizeNotif()
	case "help":
		m.io.Write(
			"command should be: clear | exit | expand | help | next | previous | remove\n",
		)
	default:
		m.io.Writef("%s not matched with any valid commands: see `help`.\n", cmd)
	}

	return false
}

func (m *Model) summarizeNotif() {
	m.queue.CallOnCurrent(func(n *Notification) {
		m.io.Write(m.conf.SummaryFormatter(n))
	})
}

func (m *Model) expandNotif() {
	m.queue.CallOnCurrent(func(n *Notification) {
		m.io.Write(m.conf.ExpandedFormatter(n))
	})
}

func (l *listener) GetServerInformation() (
	string, string, string, string, *dbus.Error,
) {
	return "painted", "none", "v0.1.3", "v1.2", nil
}

func (l *listener) GetCapabilities() ([]string, *dbus.Error) {
	return []string{"actions", "body", "persistence"}, nil
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
		Actions:   actionsFromDbus(actions),
	}

	if notif.Id == 0 { // test if we need to give this an ID
		notif.Id = atomic.AddUint32(&l.nextId, 1)
	}

	l.Recieve(notif)

	return notif.Id, nil
}

// See notification spec v1.2 table 6, row 6
func actionsFromDbus(actions []interface{}) map[string]string {
	sr_actions := make([]string, 0, len(actions))
	id_actions := make([]string, 0, len(actions))
	mp_actions := make(map[string]string)

	for i := 0; i < len(actions); i += 2 {
		id_actions = append(id_actions, actions[i].(string))
	}

	for i := 1; i < len(actions); i += 2 {
		sr_actions = append(sr_actions, actions[i].(string))
	}

	for i, _ := range sr_actions {
		mp_actions[sr_actions[i]] = id_actions[i]
	}

	return mp_actions
}

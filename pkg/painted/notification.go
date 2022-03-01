package painted

import (
	"fmt"
	"log"
	"strings"

	"github.com/lincolnauster/painted/pkg/dbus"
)

type Formatter func(*Notification) string

// This is an in-memory representation of the notification for manipulation onto
// IO. It is *not* a direct mapping of the notification spec[0] and contains
// only the information that is used by painted.
//
// [0]: https://developer-old.gnome.org/notification-spec/
type Notification struct {
	OriginApp string
	Summary   string
	Body      string
	Id        uint32
	Actions   map[string]string
}

func (n *Notification) Dismiss() {
	if n != nil {
		dbus.Emit(
			"/org/freedesktop/notifications",
			"org.freedesktop.Notifications.NotificationClosed",
			n.Id,
			uint32(2),
		)
	}
}

func (n *Notification) StringActions() string {
	action_names := make([]string, 0, len(n.Actions))
	for k, _ := range n.Actions {
		action_names = append(action_names, k)
	}
	return fmt.Sprintf("%v", action_names)
}

func (n *Notification) Format(f string) string {
	var buf strings.Builder
	var nf strings.Builder

	var state int

	// we emulate a finite state machine here. The initial state appends any
	// character to the buffer, except when it encounters a percent sign (%).
	// Then, it writes and resets the buffer it's built and sets the state
	// to 1. There, we match the next character, i.e., %o, %s, ..., and
	// either add metadata about the notification or ignore the sequence.
	for _, c := range f {
		if state == 0 && c == '%' {
			nf.WriteString(buf.String())
			buf.Reset()
			state = 1
		} else if state == 0 {
			buf.WriteRune(c)
		} else if state == 1 {
			switch c {
			case 'a':
				nf.WriteString(n.StringActions())
			case 'b':
				nf.WriteString(n.Body)
			case 'o':
				nf.WriteString(n.OriginApp)
			case 's':
				nf.WriteString(n.Summary)
			case 'i':
				nf.WriteString(fmt.Sprintf("%d", n.Id))
			default:
				// we encountered an unknown format control
				// char; render it literally and log that it was
				// unknown.
				nf.WriteRune('%')
				nf.WriteRune(c)
				log.Printf("Encountered unknown format character %c", c)
			}

			state = 0
		}
	}

	nf.WriteRune('\n')
	return nf.String()
}

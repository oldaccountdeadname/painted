package painted

import (
	"fmt"
	"strings"
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
			case 'b':
				nf.WriteString(n.Body)
			case 'o':
				nf.WriteString(n.OriginApp)
			case 's':
				nf.WriteString(n.Summary)
			case 'i':
				nf.WriteString(fmt.Sprintf("%d", n.Id))
			}

			state = 0
		}
	}

	nf.WriteRune('\n')

	return nf.String()
}

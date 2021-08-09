# painted
plain text notification daemon

## Overview
painted is a dead simple notification daemon: it reads commands from a file
(which may be a UNIX socket for IPC or any text file such as stdin for
simplicity), notifications from dbus, and writes output to a given file (usually
stdout). It aims to be a UNIX-ish replacement for standard notification daemons,
primarily for usage with bars.

Despite the name (a rough contraction of the initialism PTND, for Plain Text
Notification Daemon), the scope of PTND is not limited to plain text: planned
features include actions, sounds, and more. The goal is to achieve 'minimalism'
through an unrestrictive design that facilitates scripting, run-time
modification, and, of course, completeness.

## Usage

### In General
With no other notification daemons running, run the server (with the command
`painted`). If you send a notification with `notify-send`, you should see it
show up in the terminal window where you've run painted.

With no arguments, painted defaults to reading input commands from stdin.
Available commands (those currently implemented shown checked) are:

- [ ] help (list commands)
- [x] clear (print a blank line to hide the notification)
- [ ] dismiss
- [ ] remove (remove the notification from history: does not imply dismissal)
- [ ] prev (show the previous notification)
- [ ] next (show the next notification)
- [ ] action \<N> (select the n\'th action)
- [x] exit (close the server)

You may specify the input and output file anywhere, including /dev/stdin
(default in) and /dev/stdout (default out). If you want to control it from a
different file, specify it. For instance, in bars, you wouldn't want to have
input be given over stdin because there is no (convenient) way to access a
processes' stdin without an interactive shell.

If the input file ends in `.sock`, painted connects to it as a client on a UNIX
socket.

### In Polybar
Polybar provides the script input, where a running command's stdout will be
displayed as text. You can define a module with the following options, and
invoke it in a bar definition:

```ini
[module/painted]
type = custom/script
exec = "painted --input /tmp/painted.in"
tail = true
click-left = "echo 'clear' >> /tmp/painted.in"
```

### Somewhere else?

If you have a specific bar or application that would benefit from unique
usage instructions, open an issue and ask for them, or write it out
yourself and submit a pull request!

# Roadmap

There aren\'t any super specific goals for this project: I personally
don\'t need much out of a notification daemon, so if you need something
this doesn\'t have, go ahead and open an issue or implement it yourself
and send a PR. Just note that it *is* plain text, so it won\'t ever
support things like images or formatting.

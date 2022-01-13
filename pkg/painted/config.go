package painted

import (
	"github.com/BurntSushi/toml"
	"os"
)

type Config struct {
	SummaryFormatter  Formatter
	ExpandedFormatter Formatter
}

type configRaw struct {
	Formats notifFormats
}

type notifFormats struct {
	Summary  string
	Expanded string
}

func MakeConfigFromFile(path string) (Config, error) {
	conf := configRaw{
		notifFormats{
			`[%o] %s`,
			`%b | %a`,
		},
	}

	// Only read the configuration file if it's accessible.
	var toml_err error
	if _, err := os.Stat(path); err == nil {
		_, toml_err = toml.DecodeFile(path, &conf)
	}

	return Config{
		func(n *Notification) string {
			return n.Format(conf.Formats.Summary)
		},
		func(n *Notification) string {
			return n.Format(conf.Formats.Expanded)
		},
	}, toml_err
}

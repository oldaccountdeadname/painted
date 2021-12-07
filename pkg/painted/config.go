package painted

import "github.com/BurntSushi/toml"

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

	_, err := toml.DecodeFile(path, &conf)
	return Config{
		func(n *Notification) string {
			return n.Format(conf.Formats.Summary)
		},
		func(n *Notification) string {
			return n.Format(conf.Formats.Expanded)
		},
	}, err
}

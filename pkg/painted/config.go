package painted

import "github.com/BurntSushi/toml"

type Config struct {
	SummaryFormatter Formatter
}

type configRaw struct {
	NotifSummaryFormat string `toml:"notif_format"`
}

func MakeConfigFromFile(path string) (Config, error) {
	conf := configRaw{
		`[%o] %s`,
	}

	_, err := toml.DecodeFile(path, &conf)
	return Config{func(n *Notification) string {
		return n.Format(conf.NotifSummaryFormat)
	}}, err
}

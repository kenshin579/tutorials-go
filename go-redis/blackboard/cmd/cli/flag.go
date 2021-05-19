package cli

import "flag"

type Flags struct {
	ConfigPath string
}

func ParseFlags() *Flags {
	flags := &Flags{}
	flag.StringVar(&flags.ConfigPath, "config-file", "config/config.yaml", "path to config file")
	flag.Parse()

	return flags
}

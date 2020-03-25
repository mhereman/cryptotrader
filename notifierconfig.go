package cryptotrader

import "strings"

// NotifierConfig represents the config for the notifier to use
type NotifierConfig struct {
	// Name of the notifier to use
	Name string

	// ArgMap arguments of the notifier
	ArgMap map[string]string
}

// NewNotifierConfigFromFlags creates a new NotifierConfig insance from the cmdline argument values
func NewNotifierConfigFromFlags(name string, args map[string]string) (nc NotifierConfig, err error) {
	nc.Name = strings.ToLower(name)
	nc.ArgMap = args
	return
}

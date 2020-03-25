package cryptotrader

import "strings"

// ExchangeConfig represents the config for the exchange to trade on
type ExchangeConfig struct {
	// Name of the exchange
	Name string

	// Exchange specific map of arguments
	ArgMap map[string]string
}

// NewExchangeConfigFromFlags creates a new ExchangeConfig insance from the cmdline argument values
func NewExchangeConfigFromFlags(name string, args map[string]string) (ec ExchangeConfig, err error) {
	ec.Name = strings.ToLower(name)
	ec.ArgMap = args
	return
}

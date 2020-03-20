package exchange

import (
	"context"
	"fmt"

	"github.com/mhereman/cryptotrader/interfaces"
	"github.com/mhereman/cryptotrader/logger"
)

var exchangeFactory map[string]func(context.Context, map[string]string) (interfaces.IExchangeDriver, error) = make(map[string]func(context.Context, map[string]string) (interfaces.IExchangeDriver, error))

func RegisterExchange(name string, factory func(context.Context, map[string]string) (interfaces.IExchangeDriver, error)) {
	exchangeFactory[name] = factory
	logger.Printf("Registered exchanged: %s\n", name)
}

func GetExchange(ctx context.Context, name string, args map[string]string) (exchange interfaces.IExchangeDriver, err error) {
	var ok bool
	var fn func(context.Context, map[string]string) (interfaces.IExchangeDriver, error)

	if fn, ok = exchangeFactory[name]; !ok {
		err = fmt.Errorf("Exchange %s does not exist")
		return
	}

	exchange, err = fn(ctx, args)
	return
}

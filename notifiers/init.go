package notifiers

import (
	"context"
	"fmt"

	"github.com/mhereman/cryptotrader/interfaces"
	"github.com/mhereman/cryptotrader/logger"
)

var notifierFactory map[string]func(context.Context, map[string]string) (interfaces.INotifier, error) = make(map[string]func(context.Context, map[string]string) (interfaces.INotifier, error))

func RegisterNotifier(name string, factory func(context.Context, map[string]string) (interfaces.INotifier, error)) {
	notifierFactory[name] = factory
	logger.Printf("Registered notifier: %s\n", name)
}

func GetNotifier(ctx context.Context, name string, args map[string]string) (notifier interfaces.INotifier, err error) {
	var ok bool
	var fn func(context.Context, map[string]string) (interfaces.INotifier, error)

	if fn, ok = notifierFactory[name]; !ok {
		err = fmt.Errorf("Notifier %s does not exist")
		return
	}

	notifier, err = fn(ctx, args)
	return
}

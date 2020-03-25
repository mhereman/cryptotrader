package notifiers

import (
	"context"
	"fmt"

	"github.com/mhereman/cryptotrader/interfaces"
	"github.com/mhereman/cryptotrader/logger"
)

var notifierFactory map[string]func(context.Context, map[string]string) (interfaces.INotifier, error) = make(map[string]func(context.Context, map[string]string) (interfaces.INotifier, error))

// RegisterNotifier registers a new notifier factory function.
// This functio nshould be called from the init() function of the
// notifier plugin.
func RegisterNotifier(name string, factory func(context.Context, map[string]string) (interfaces.INotifier, error)) {
	notifierFactory[name] = factory
	logger.Printf("Registered notifier: %s\n", name)
}

// GetNotifier returns a notifier plugin registered under the provided name or an error
func GetNotifier(ctx context.Context, name string, args map[string]string) (notifier interfaces.INotifier, err error) {
	var ok bool
	var fn func(context.Context, map[string]string) (interfaces.INotifier, error)

	if fn, ok = notifierFactory[name]; !ok {
		err = fmt.Errorf("Notifier %s does not exist", name)
		return
	}

	notifier, err = fn(ctx, args)
	return
}

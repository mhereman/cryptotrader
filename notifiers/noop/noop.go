package noop

import (
	"context"

	"github.com/mhereman/cryptotrader/interfaces"
	"github.com/mhereman/cryptotrader/notifiers"
)

const notifierName = "noop"

func init() {
	notifiers.RegisterNotifier(notifierName, createNoop)
}

// Noop represents the NO-OP notifier
// This notifier does not do anything
type Noop struct {
}

func createNoop(ctx context.Context, config map[string]string) (notifier interfaces.INotifier, err error) {
	notifier = &Noop{}
	return
}

// Name returns the name of the notifier
func (noop Noop) Name() string {
	return notifierName
}

// Notify ...
func (noop Noop) Notify(context.Context, []byte) (err error) {
	return
}

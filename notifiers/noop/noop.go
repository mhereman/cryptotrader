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

type Noop struct {
}

func createNoop(ctx context.Context, config map[string]string) (notifier interfaces.INotifier, err error) {
	notifier = &Noop{}
	return
}

func (noop Noop) Name() string {
	return notifierName
}

func (noop Noop) Notify(context.Context, []byte) (err error) {
	return
}

package interfaces

import "context"

// INotifier represents the notifier plugin interface
type INotifier interface {
	// Name returs the name of the notifier plugin
	Name() string

	// Notify executs the notifier with the provided message
	Notify(context.Context, []byte) error
}

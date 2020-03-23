package interfaces

import "context"

type INotifier interface {
	Name() string

	Notify(context.Context, []byte) error
}

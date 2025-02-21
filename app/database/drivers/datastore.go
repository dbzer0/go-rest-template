package drivers

import "context"

type DataStore interface {
	Name() string
	Close(ctx context.Context) error
	Connect() error

	// расширяем функционал datastore здесь
}

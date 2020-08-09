package drivers

import "context"

type DataStore interface {
	// расширяем функционал datastore здесь
	Name() string
	Close(ctx context.Context) error
	Connect() error
}

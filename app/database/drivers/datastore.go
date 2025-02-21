package drivers

import "context"

type DataStore interface {
	Name() string
	Close(ctx context.Context) error
	Connect() error

	// расширяем функционал DataStore здесь
	OrdersRepository
}

type OrdersRepository interface {
	Orders() error
}

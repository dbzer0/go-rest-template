package database

import (
	"fmt"
	"sync"

	"github.com/dbzer0/go-rest-template/app/database/drivers"
	"github.com/dbzer0/go-rest-template/app/database/drivers/mongo"
	"github.com/pkg/errors"
)

// DataStoreFactory определяет функцию создания DataStore
type DataStoreFactory func(conf drivers.DataStoreConfig) (drivers.DataStore, error)

// Registry управляет регистрацией фабрик датасторов
type Registry struct {
	mu        sync.RWMutex
	factories map[string]DataStoreFactory
}

// NewRegistry создает новый реестр датасторов
func NewRegistry() *Registry {
	r := &Registry{
		factories: make(map[string]DataStoreFactory),
	}
	r.registerDefaults()
	return r
}

// registerDefaults регистрирует датасторы по умолчанию
func (r *Registry) registerDefaults() {
	// Здесь можно добавить больше драйверов
	must(r.Register("mongo", mongo.New))
}

// Register регистрирует новую фабрику датастора
func (r *Registry) Register(name string, factory DataStoreFactory) error {
	if factory == nil {
		return fmt.Errorf("datastore factory '%s' cannot be nil", name)
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.factories[name]; exists {
		return fmt.Errorf("datastore factory '%s' is already registered", name)
	}

	r.factories[name] = factory
	return nil
}

// Create создает новый экземпляр датастора
func (r *Registry) Create(conf drivers.DataStoreConfig) (drivers.DataStore, error) {
	r.mu.RLock()
	factory, exists := r.factories[conf.DataStoreName]
	r.mu.RUnlock()

	if !exists {
		return nil, NewErrInvalidDataStore(r.AvailableDrivers())
	}

	return factory(conf)
}

// AvailableDrivers возвращает список доступных драйверов
func (r *Registry) AvailableDrivers() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	drivers := make([]string, 0, len(r.factories))
	for name := range r.factories {
		drivers = append(drivers, name)
	}
	return drivers
}

// DefaultRegistry глобальный реестр по умолчанию
var DefaultRegistry = NewRegistry()

// Connect создает и подключает датастор
func Connect(conf drivers.DataStoreConfig) (drivers.DataStore, error) {
	ds, err := DefaultRegistry.Create(conf)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create datastore")
	}

	if err = ds.Connect(); err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	return ds, nil
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// ErrInvalidDataStore представляет ошибку неверного датастора
type ErrInvalidDataStore struct {
	AvailableDrivers []string
}

func NewErrInvalidDataStore(drivers []string) error {
	return &ErrInvalidDataStore{AvailableDrivers: drivers}
}

func (e *ErrInvalidDataStore) Error() string {
	return fmt.Sprintf("invalid datastore name, available drivers: %v", e.AvailableDrivers)
}

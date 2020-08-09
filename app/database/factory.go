package database

import (
	"log"

	"github.com/dbzer0/go-rest-template/app/database/drivers"
	"github.com/dbzer0/go-rest-template/app/database/drivers/mongo"

	"github.com/pkg/errors"
)

type DataStoreFactory func(conf drivers.DataStoreConfig) (drivers.DataStore, error)

// dataStoreFactories содержит все доступыне DataStore фабрики
var dataStoreFactories = make(map[string]DataStoreFactory)

// Register регистрирует новую фабрику.
// Если фабрика была зарегистрирована раньше, то сообщает об этом.
// Если фабрика не существует, то сообщает об этом.
func Register(name string, factory DataStoreFactory) error {
	if factory == nil {
		return errors.Errorf("datastore factory %s does not exist.", name)
	}

	_, registered := dataStoreFactories[name]
	if registered {
		log.Printf("[WARNING] datastore factory %s already registered. Ignoring.", name)
	} else {
		dataStoreFactories[name] = factory
	}

	return nil
}

func init() {
	if err := Register("mongo", mongo.New); err != nil {
		log.Panic(err)
	}
}

func New(conf drivers.DataStoreConfig) (drivers.DataStore, error) {
	engineFactory, ok := dataStoreFactories[conf.DataStoreName]
	if !ok {
		// выбран неверный datastore, получаем список доступных и отдаем его
		// пользователю
		availableDataStores := make([]string, 0, len(dataStoreFactories)-1)
		for k := range dataStoreFactories {
			availableDataStores = append(availableDataStores, k)
		}

		return nil, ErrInvalidDataStoreName.Error(availableDataStores)
	}

	return engineFactory(conf)
}

// Connect подключение к БД
func Connect(conf drivers.DataStoreConfig) (drivers.DataStore, error) {
	ds, err := New(conf)
	if err != nil {
		return nil, errors.Wrap(err, "cannot create datastore")
	}

	if err = ds.Connect(); err != nil {
		return nil, errors.Wrap(err, "cannot connect to database")
	}

	return ds, nil
}

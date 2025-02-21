package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	"github.com/dbzer0/go-rest-template/app/database/drivers"
)

const (
	connectionTimeout = 3 * time.Second
	ensureIdxTimeout  = 10 * time.Second
)

type Mongo struct {
	connURL string
	dbName  string

	client *mongo.Client
	db     *mongo.Database

	connectionTimeout time.Duration
	ensureIdxTimeout  time.Duration

	// TODO: Встраиваем интерфейсы репозиториев
	drivers.OrdersRepository
}

func New(conf drivers.DataStoreConfig) (drivers.DataStore, error) {
	if conf.URL == "" {
		return nil, drivers.ErrInvalidConfigStruct
	}

	if conf.DataBaseName == "" {
		return nil, drivers.ErrInvalidConfigStruct
	}

	return &Mongo{
		connURL:           conf.URL,
		dbName:            conf.DataBaseName,
		connectionTimeout: connectionTimeout,
		ensureIdxTimeout:  ensureIdxTimeout,
	}, nil
}

func (m *Mongo) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	var err error
	m.client, err = mongo.Connect(ctx, options.Client().ApplyURI(m.connURL))
	if err != nil {
		return err
	}

	if err = m.Ping(); err != nil {
		return err
	}

	m.db = m.client.Database(m.dbName)

	// TODO: Инициализируем репозитории и встраиваем его
	m.OrdersRepository = NewOrdersRepository(m.db)

	// Автоматически создаем индексы
	return m.ensureIndexes()
}

func (m *Mongo) ensureIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.ensureIdxTimeout)
	defer cancel()

	// TODO: Собираем все репозитории, отвечающие за создание индексов.
	ensurers := []IndexEnsurer{
		NewOrdersRepository(m.db),
	}

	for _, e := range ensurers {
		if err := e.EnsureIndexes(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (m *Mongo) Name() string { return "Mongo" }

func (m *Mongo) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	return m.client.Ping(ctx, readpref.Primary())
}

func (m *Mongo) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

// Вспомогательная функция для получения существующих индексов
func (m *Mongo) existingIndexes(ctx context.Context, col *mongo.Collection) (map[string]bool, error) {
	indexes := make(map[string]bool)
	cursor, err := col.Indexes().List(ctx)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		indexName := cursor.Current.Lookup("name").StringValue()
		indexes[indexName] = true
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return indexes, nil
}

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
	DB     *mongo.Database

	connectionTimeout time.Duration
	ensureIdxTimeout  time.Duration
}

func (m *Mongo) Name() string { return "Mongo" }

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
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	m.client, err = mongo.Connect(ctx, options.Client().ApplyURI(m.connURL))
	if err != nil {
		return err
	}

	if err := m.Ping(); err != nil {
		return err
	}

	m.DB = m.client.Database(m.dbName)

	// убеждаемся что созданы все необходимые индексы
	return m.ensureIndexes()
}

func (m *Mongo) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	return m.client.Ping(ctx, readpref.Primary())
}

func (m *Mongo) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}

// убеждается что все индексы построены
func (m *Mongo) ensureIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.connectionTimeout)
	defer cancel()

	// TODO
	_ = ctx

	return nil
}

// indexExistsByName проверяет существование индекса с именем name.
func (m *Mongo) indexExistsByName(ctx context.Context, collection *mongo.Collection, name string) (bool, error) {
	cur, err := collection.Indexes().List(ctx)
	if err != nil {
		return false, err
	}

	for cur.Next(ctx) {
		if name == cur.Current.Lookup("name").StringValue() {
			return true, nil
		}
	}

	return false, nil
}

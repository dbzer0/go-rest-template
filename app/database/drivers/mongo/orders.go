package mongo

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type OrdersRepository struct {
	collection       *mongo.Collection
	ensureIdxTimeout time.Duration
}

func NewOrdersRepository(db *mongo.Database) *OrdersRepository {
	return &OrdersRepository{
		collection:       db.Collection("orders"),
		ensureIdxTimeout: ensureIdxTimeout,
	}
}

func (o *OrdersRepository) Orders() error {
	fmt.Println("creating orders example")
	return nil
}

func (o *OrdersRepository) EnsureIndexes(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, o.ensureIdxTimeout)
	defer cancel()

	//idxs := []indexSpec{
	//	// Индекс по client_id
	//	{"client_id_1", bson.D{{"client_id", 1}}, false},
	//	// Индекс по timestamp
	//	{"timestamp_1", bson.D{{"timestamp", 1}}, false},
	//	// Составной индекс по client_id и timestamp
	//	{"client_id_timestamp", bson.D{{"client_id", 1}, {"timestamp", 1}}, false},
	//}

	idxs := []indexSpec{}
	fmt.Println("ensuring indexes for", o.collection.Name(), "example")
	return ensureIndexesForCollection(ctx, o.collection, idxs, o.ensureIdxTimeout)
}

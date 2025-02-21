package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IndexEnsurer interface {
	EnsureIndexes(ctx context.Context) error
}

// indexSpec определяет спецификацию для создания индекса.
type indexSpec struct {
	Name   string
	Keys   bson.D
	Unique bool
}

// ensureIndexesForCollection проверяет и создаёт отсутствующие индексы для заданной коллекции.
func (m *Mongo) ensureIndexesForCollection(ctx context.Context, col *mongo.Collection, idxs []indexSpec, maxTime time.Duration) error {
	existingIndexes, err := m.existingIndexes(ctx, col)
	if err != nil {
		return err
	}

	var models []mongo.IndexModel
	for _, index := range idxs {
		if _, exists := existingIndexes[index.Name]; !exists {
			opts := options.Index().SetName(index.Name)
			if index.Unique {
				opts.SetUnique(true)
			}
			models = append(models, mongo.IndexModel{
				Keys:    index.Keys,
				Options: opts,
			})
		}
	}

	if len(models) == 0 {
		return nil // Нет новых индексов для создания.
	}

	createOpts := options.CreateIndexes().SetMaxTime(maxTime)
	_, err = col.Indexes().CreateMany(ctx, models, createOpts)
	return err
}

// Вспомогательная функция для создания индексов для конкретной коллекции.
func ensureIndexesForCollection(ctx context.Context, col *mongo.Collection, idxs []indexSpec, maxTime time.Duration) error {
	existingIndexes, err := existingIndexes(ctx, col)
	if err != nil {
		return err
	}

	var models []mongo.IndexModel
	for _, index := range idxs {
		if _, exists := existingIndexes[index.Name]; !exists {
			opts := options.Index().SetName(index.Name)
			if index.Unique {
				opts.SetUnique(true)
			}
			models = append(models, mongo.IndexModel{
				Keys:    index.Keys,
				Options: opts,
			})
		}
	}

	if len(models) == 0 {
		return nil // Нет новых индексов для создания.
	}

	createOpts := options.CreateIndexes().SetMaxTime(maxTime)
	_, err = col.Indexes().CreateMany(ctx, models, createOpts)
	return err
}

// Функция для получения существующих индексов.
func existingIndexes(ctx context.Context, col *mongo.Collection) (map[string]bool, error) {
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
	return indexes, cursor.Err()
}

package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// indexSpec определяет спецификацию для создания индекса.
type indexSpec struct {
	Name   string
	Keys   bson.D
	Unique bool
}

func (m *Mongo) ensureIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.ensureIdxTimeout)
	defer cancel()

	//if err := m.ensureAPIUsageIndexes(ctx); err != nil {
	//	return err
	//}
	_ = ctx

	return nil
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

//func (m *Mongo) ensureAPIUsageIndexes(ctx context.Context) error {
//	col := m.db.Collection(collectionAPIUsage)
//	idxs := []indexSpec{
//		// Индекс по client_id.
//		// Позволяет быстро находить все записи использования API для заданного клиента.
//		{"client_id_1", bson.D{{"client_id", 1}}, false},
//
//		// Индекс по timestamp.
//		// Обеспечивает быстрый поиск записей по временной метке вызова API (например, для выборок по диапазону дат).
//		{"timestamp_1", bson.D{{"timestamp", 1}}, false},
//
//		// Составной индекс по client_id и timestamp.
//		// Ускоряет выполнение запросов, фильтрующих данные по клиенту и по диапазону дат одновременно.
//		{"client_id_timestamp", bson.D{{"client_id", 1}, {"timestamp", 1}}, false},
//	}
//	// Здесь в качестве maxTime используем 5 секунд.
//	return m.ensureIndexesForCollection(ctx, col, idxs, 5*time.Second)
//}

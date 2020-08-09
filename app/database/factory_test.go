package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// При добавлении новых фабрик, это число нужно увеличивать
var fabricCount = 1

// Тест на проверку регистрации фабрик. Проверяет их наличие и адресацию фабрик.
func TestRegister(t *testing.T) {
	// Проверяет количество доступных датасторов
	assert.Equal(t, len(dataStoreFactories), fabricCount)

	// в существующем датасторе обязательно должен быть адрес фабрики
	assert.NotNil(t, dataStoreFactories["mongo"])

	// у неправильного датастора будет адрес фабрики nil
	assert.Nil(t, dataStoreFactories["invalid_datastore"])
}

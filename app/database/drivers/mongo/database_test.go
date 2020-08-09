package mongo

import (
	"context"
	"os"
	"testing"

	"github.com/dbzer0/go-rest-template/app/database/drivers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// TestDatabaseSuite запускает suit для тестирования БД.
func TestDatabaseSuite(t *testing.T) {
	suite.Run(t, new(DatabaseSuite))
}

// DatabaseSuite содержит suit для тестирования БД.
type DatabaseSuite struct {
	suite.Suite
	db drivers.DataStore
}

// BeforeTest формирует конфигурацию для MongoDB и пробует
// присоединиться к БД.
func (s *DatabaseSuite) BeforeTest(string, string) {
	url := os.Getenv("MONGO_URL")
	if url == "" {
		url = "mongodb://localhost:27017"
	}

	conf := drivers.DataStoreConfig{
		URL:           url,
		DataStoreName: "mongo",
		DataBaseName:  "test-db",
	}
	db, _ := New(conf)

	err := db.Connect()
	assert.NoError(s.T(), err, "cannot connect to MongoDB server")
	if err != nil {
		s.T().Fatal(err)
		return
	}

	s.db = db
}

// AfterTest закрывает соединение к БД.
func (s *DatabaseSuite) AfterTest(string, string) {
	if s.db != nil {
		s.db.Close(context.Background())
	}
}

package test

import (
	"context"

	"github.com/dbzer0/go-rest-template/app/database/drivers"
)

type Manager struct {
	ds drivers.DataStore
}

// NewManager создаёт новый менеджер профилей клиентов.
func NewManager(ds drivers.DataStore) *Manager {
	return &Manager{ds: ds}
}

func (m *Manager) Create(ctx context.Context) error {
	return nil
}

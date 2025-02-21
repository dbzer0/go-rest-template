package orders

import (
	"context"

	"github.com/dbzer0/go-rest-template/app/database/drivers"
)

type Manager struct {
	repo drivers.OrdersRepository
}

// NewManager создаёт новый менеджер профилей клиентов.
func NewManager(repo drivers.OrdersRepository) *Manager {
	return &Manager{repo: repo}
}

func (m *Manager) Create(ctx context.Context) error {
	m.repo.Orders()
	return nil
}

package api

import (
	"github.com/go-chi/chi/v5"
)

// AdminAPI инкапсулирует админские API-эндпоинты.
type AdminAPI struct {
}

// NewAdminAPI создаёт новый экземпляр AdminAPI.
func NewAdminAPI() *AdminAPI {
	return &AdminAPI{}
}

// Routes возвращает маршруты для админского API.
func (a *AdminAPI) Routes() chi.Router {
	r := chi.NewRouter()

	return r
}

// StatusResponse используется для возврата статуса операции.
type StatusResponse struct {
	Status string `json:"status"`
}

package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

const (
	defaultTimeout = 10 * time.Second
)

type API struct {
}

func NewAPI() *API {
	return &API{}
}

func (a API) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/client/test", a.TestHandler)
	return r
}

func (a *API) TestHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "Access granted"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

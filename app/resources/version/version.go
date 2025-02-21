package version

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// response - ответ на запрос версии.
type response struct {
	APIVersion     string `json:"api"`
	ServiceVersion string `json:"version"`
}

func NewResource(serviceVersion, apiVersion string) *Resource {
	return &Resource{
		ServiceVersion: serviceVersion,
		APIVersion:     apiVersion,
	}
}

// Resource - структура содержащая версию APIVersion и приложения.
type Resource struct {
	ServiceVersion string
	APIVersion     string
}

func (vr Resource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", vr.Version)

	return r
}

func (vr Resource) Version(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, response{
		APIVersion:     vr.APIVersion,
		ServiceVersion: vr.ServiceVersion,
	})
}

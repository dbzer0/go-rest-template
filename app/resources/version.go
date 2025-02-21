package resources

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const APIVersion = "v1"

// VersionResponse - ответ на запрос версии.
type VersionResponse struct {
	API     string `json:"api"`
	Version string `json:"version"`
}

func NewVersionResponse(version string) *VersionResource {
	return &VersionResource{version: version}
}

// VersionResource - структура содержащая версию API и приложения.
type VersionResource struct {
	version string
}

func (vr VersionResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", vr.Version)

	return r
}

func (vr VersionResource) Version(w http.ResponseWriter, r *http.Request) {
	render.JSON(w, r, VersionResponse{
		API:     APIVersion,
		Version: vr.version,
	})
}

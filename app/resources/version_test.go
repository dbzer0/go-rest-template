package resources

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

// Тестирует uri "/version" на валидность ответа, код возврата и тип контента
func TestGetVersion(t *testing.T) {
	req, _ := http.NewRequest("GET", "/version", nil)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Mount("/version", VersionResource{}.Routes())
	version := fmt.Sprintf("{\"api\":\"%s\",\"version\":\"\"}\n", APIVersion)

	r.ServeHTTP(w, req)

	assert.Equal(t, w.Body.String(), version, "VersionResponse should be"+version+", was "+w.Body.String())
	assert.Equal(t, w.Code, http.StatusOK, "Response HTTP code should be 200, was: %d", w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"),
		"Content-Type should be application/json, was "+w.Header().Get("Content-Type"))
}

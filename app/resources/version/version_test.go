package version

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

// Тестирует uri "/ServiceVersion" на валидность ответа, код возврата и тип контента
func TestVersion(t *testing.T) {
	req, _ := http.NewRequest(http.MethodGet, "/version", nil)
	w := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Mount("/ServiceVersion", Resource{}.Routes())
	version := fmt.Sprintf("{\"api\":\"%s\",\"ServiceVersion\":\"\"}\n", APIVersion)

	r.ServeHTTP(w, req)

	assert.Equal(t, w.Body.String(), version, "response should be"+version+", was "+w.Body.String())
	assert.Equal(t, w.Code, http.StatusOK, "response HTTP code should be 200, was: %d", w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"),
		"Content-Type should be application/json, was "+w.Header().Get("Content-Type"))
}

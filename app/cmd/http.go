package cmd

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/dbzer0/go-rest-template/app/resources/api"
	"github.com/dbzer0/go-rest-template/app/resources/version"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

const (
	compressLevel   = 5
	readTimeout     = 5 * time.Second
	writeTimeout    = 30 * time.Second
	shutdownTimeout = 10 * time.Second
)

// HTTPCommand инкапсулирует запуск HTTP-сервера.
type HTTPCommand struct {
	srv *httpServer
}

// NewHTTPCommand создаёт новый HTTPCommand.
func NewHTTPCommand(ctx context.Context, opts *Configuration, version string) *HTTPCommand {
	return &HTTPCommand{
		srv: newHTTPServer(ctx, opts, version),
	}
}

// Execute запускает HTTP-сервер.
func (c *HTTPCommand) Execute(ctx context.Context) error {
	return c.srv.Run()
}

// httpServer содержит настройки и состояние HTTP-сервера.
type httpServer struct {
	Address           string
	CertFile, KeyFile *string
	BasePath          string
	serviceVersion    string

	masterCtx context.Context
	wg        sync.WaitGroup
	IsTesting bool
}

// newHTTPServer инициализирует httpServer, заполняя его конфигурацией.
func newHTTPServer(ctx context.Context, opts *Configuration, version string) *httpServer {
	return &httpServer{
		masterCtx:      ctx,
		Address:        opts.ListenAddr,
		BasePath:       opts.BasePath,
		serviceVersion: version,
		IsTesting:      opts.IsTesting,
		CertFile:       nonEmptyPtr(opts.CertFile),
		KeyFile:        nonEmptyPtr(opts.KeyFile),
	}
}

// nonEmptyPtr возвращает указатель на строку, если она не пустая.
func nonEmptyPtr(s string) *string {
	if s != "" {
		return &s
	}
	return nil
}

// setupRouter настраивает маршруты и middleware для сервера.
func (srv *httpServer) setupRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.NoCache)
	r.Use(middleware.Compress(compressLevel))
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins(srv.IsTesting),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	const apiVersion = "v1"

	// Монтируем обработчики
	r.Mount("/version", version.NewResource(srv.serviceVersion, apiVersion).Routes())
	r.Mount("/api/v1", api.NewAPI().Routes())

	return r
}

// allowedOrigins возвращает список разрешённых источников в зависимости от режима тестирования.
func allowedOrigins(testing bool) []string {
	if testing {
		return []string{"*"}
	}
	return []string{}
}

// Run запускает HTTP-сервер и обрабатывает graceful shutdown.
func (srv *httpServer) Run() error {
	router := srv.setupRouter()

	server := &http.Server{
		Addr:         srv.Address,
		Handler:      router,
		BaseContext:  func(_ net.Listener) context.Context { return srv.masterCtx },
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}

	// Запускаем горутину для graceful shutdown
	srv.wg.Add(1)
	go func() {
		defer srv.wg.Done()
		<-srv.masterCtx.Done()
		ctxShutdown, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := server.Shutdown(ctxShutdown); err != nil {
			log.Printf("[ERROR] HTTP server shutdown failed: %v", err)
			return
		}
		log.Printf("[INFO] HTTP server stopped gracefully")
	}()

	log.Printf("[INFO] serving HTTP server on %q", srv.Address)

	var err error
	if srv.CertFile == nil && srv.KeyFile == nil {
		err = server.ListenAndServe()
	} else {
		err = server.ListenAndServeTLS(*srv.CertFile, *srv.KeyFile)
	}

	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	srv.wg.Wait() // Ждём завершения shutdown
	return nil
}

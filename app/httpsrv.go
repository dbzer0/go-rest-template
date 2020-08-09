package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dbzer0/go-rest-template/app/database/drivers"
	"github.com/dbzer0/go-rest-template/app/resources"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

const compressLevel = 5

type HTTPServer struct {
	Address           string
	CertFile, KeyFile *string
	ds                drivers.DataStore
	BasePath          string
	version           string // версия приложения
	masterCtx         context.Context
	idleConnsClosed   chan struct{}
	IsTesting         bool
}

func NewHTTPServer(ctx context.Context, opts *configuration, ds drivers.DataStore, version string) *HTTPServer {
	srv := &HTTPServer{
		masterCtx:       ctx,
		Address:         opts.ListenAddr,
		BasePath:        opts.BasePath,
		ds:              ds,
		version:         version,             // версия приложения
		idleConnsClosed: make(chan struct{}), // способ определить незавершенные соединения
		IsTesting:       opts.IsTesting,
	}

	if opts.CertFile != "" {
		srv.CertFile = &opts.CertFile
	}

	if opts.KeyFile != "" {
		srv.KeyFile = &opts.KeyFile
	}

	return srv
}

// setupRouter инициализирует HTTP роутер.
// Функция используется для подключения middleware и маппинга ресурсов.
func (srv *HTTPServer) setupRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.NoCache) // no-cache
	//r.Use(middleware.RequestID) // вставляет request ID в контекст каждого запроса
	r.Use(middleware.Logger)    // логирует начало и окончание каждого запроса с указанием времени обработки
	r.Use(middleware.Recoverer) // управляемо обрабатывает паники и выдает stack trace при их возникновении
	r.Use(middleware.RealIP)    // устанавливает RemoteAddr для каждого запроса с заголовками X-Forwarded-For или X-Real-IP
	r.Use(middleware.NewCompressor(compressLevel).Handler)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins(srv.IsTesting),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// монтируем дополнительные ресурсы
	r.Mount("/version", resources.VersionResource{Version: srv.version}.Routes())

	return r
}

// getAllowedOrigins возвращает список хостов для C.O.R.S.
func allowedOrigins(testing bool) []string {
	if testing {
		return []string{"*"}
	}

	return []string{}
}

// Run запускает HTTP или HTTPS листенер в зависимости от того как заполнена
// структура HTTPServer{}.
func (srv *HTTPServer) Run() error {
	const (
		readTimeout  = 5 * time.Second
		writeTimeout = 30 * time.Second
	)

	s := http.Server{
		Addr:         srv.Address,
		Handler:      chi.ServerBaseContext(srv.masterCtx, srv.setupRouter()),
		ReadTimeout:  readTimeout,  // wait() + tls handshake + req.headers + req.body
		WriteTimeout: writeTimeout, // все что выше + response
	}

	go func() {
		<-srv.masterCtx.Done()

		if err := s.Shutdown(srv.masterCtx); err != nil {
			log.Printf("[ERROR] HTTP server Shutdown: %v", err)
		}

		srv.Wait()
	}()

	log.Printf("[INFO] serving HTTP on \"%s\"", srv.Address)

	var err error
	if srv.CertFile == nil && srv.KeyFile == nil {
		if err = s.ListenAndServe(); err != nil {
			close(srv.idleConnsClosed)
			return err
		}
	} else {
		if err = s.ListenAndServeTLS(*srv.CertFile, *srv.KeyFile); err != nil {
			close(srv.idleConnsClosed)
			return err
		}
	}

	return nil
}

// Wait ожидает момента завершения обработки всех соединений.
func (srv *HTTPServer) Wait() {
	<-srv.idleConnsClosed
}

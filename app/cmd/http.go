package cmd

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/dbzer0/go-rest-template/app/resources"
	"github.com/dbzer0/go-rest-template/app/resources/api"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

const compressLevel = 5

type HTTPCommand struct {
	srv *httpServer
}

func NewHTTPCommand(
	ctx context.Context,
	opts *Configuration,
	version string,
) *HTTPCommand {
	return &HTTPCommand{
		srv: newHTTPServer(ctx, opts, version),
	}
}

func (c *HTTPCommand) Execute(ctx context.Context) error {
	return c.srv.Run()
}

type httpServer struct {
	Address           string
	CertFile, KeyFile *string
	BasePath          string
	version           string
	masterCtx         context.Context
	idleConnsClosed   chan struct{}
	IsTesting         bool
}

func newHTTPServer(ctx context.Context, opts *Configuration, version string) *httpServer {
	srv := &httpServer{
		masterCtx:       ctx,
		Address:         opts.ListenAddr,
		BasePath:        opts.BasePath,
		version:         version,
		idleConnsClosed: make(chan struct{}),
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

func (srv *httpServer) setupRouter() chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.NoCache)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.NewCompressor(compressLevel).Handler)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins(srv.IsTesting),
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Mount("/version", resources.NewVersionResponse(srv.version).Routes())
	r.Mount("/api/v1", api.NewAPI().Routes())

	return r
}

func allowedOrigins(testing bool) []string {
	if testing {
		return []string{"*"}
	}
	return []string{}
}

func (srv *httpServer) Run() error {
	const (
		readTimeout  = 5 * time.Second
		writeTimeout = 30 * time.Second
	)

	s := http.Server{
		Addr:         srv.Address,
		Handler:      chi.ServerBaseContext(srv.masterCtx, srv.setupRouter()),
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
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

func (srv *httpServer) Wait() {
	<-srv.idleConnsClosed
}

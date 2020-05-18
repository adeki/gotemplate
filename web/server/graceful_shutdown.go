package server

import (
  "net/http"
  "github.com/go-chi/chi"
)

type Server struct {}

func New() *Server {
  return &Server{}
}

func (s *Server) Handler() http.Handler {
  // r := chi.NewRouter()
  // return r
  return nil
}

func (s *Server) Run() {
	srv := http.Server{
		Addr:    ":8080",
		Handler: s.Handler(),
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		// We received sigint or sigterm, shut down.
		logger.Info("shutdown server...")
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			logger.Warnf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	logger.Info("starting server...")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		logger.Warnf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

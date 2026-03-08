package server

import (
	"net/http"
	"time"

	"fg_bot/internal/server/handler"
)

type Server struct {
	srv *http.Server
}

func New(h *handler.Handler, port string) *Server {
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      h,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{srv: srv}
}

func (s *Server) ListenAndServe() error {
	return s.srv.ListenAndServe()
}

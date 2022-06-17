package internalserver

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"tec-doc/internal/web"
)

type Server struct {
	router *gin.Engine
	server http.Server
}

func NewInternalServer(bindingAddress string) web.Server {
	router := initRouter()
	return &Server{
		router: router,
		server: http.Server{
			Addr:    bindingAddress,
			Handler: router,
		},
	}
}

func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.server.Shutdown(context.Background())
}

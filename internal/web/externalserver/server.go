package externalserver

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"tec-doc/internal/web"
)

type externalHttpServer struct {
	router *gin.Engine
	server http.Server
}

func NewExternalServer(bindingAddress string) web.Server {
	router := initExternalRouter()
	return &externalHttpServer{
		router: router,
		server: http.Server{
			Addr:    bindingAddress,
			Handler: router,
		},
	}
}

func (s *externalHttpServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *externalHttpServer) Stop() error {
	return s.server.Shutdown(context.Background())
}

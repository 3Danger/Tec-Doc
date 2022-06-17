package internalserver

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"tec-doc/internal/web"
)

type internalHttpServer struct {
	router  *gin.Engine
	server  *http.Server
	metrics *Metrics
}

func NewInternalServer(bindingAddress string) web.Server {
	router := initInternalRouter()
	return &internalHttpServer{
		router:  router,
		metrics: NewMetric("Server", "Internal"),
		server: &http.Server{
			Addr:    bindingAddress,
			Handler: router,
		},
	}
}

func (s *internalHttpServer) Start() error {
	return s.server.ListenAndServe()
}

func (s *internalHttpServer) Stop() error {
	return s.server.Shutdown(context.Background())
}

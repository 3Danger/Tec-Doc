package internalserver

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"tec-doc/pkg/ginLogger"
)

type Server interface {
	Start() error
	Stop() error
}

type internalHttpServer struct {
	router *gin.Engine
	server *http.Server
}

func New(bindingPort string) *internalHttpServer {
	serv := new(internalHttpServer)

	serv.router = gin.New()
	serv.configureRouter()
	serv.server = &http.Server{
		Addr:    bindingPort,
		Handler: serv.router,
	}
	return serv
}

func (i *internalHttpServer) configureRouter() {
	i.router.Use(gin.Recovery())
	i.router.Use(ginLogger.Logger(os.Stdout))
	i.router.GET("/health", i.Health)
	i.router.GET("/readiness", i.Readiness)
	i.router.GET("/metrics", i.Metrics)
}

func (i *internalHttpServer) Start() error {
	log.Info().Msg("start internal server on " + i.server.Addr)
	return i.server.ListenAndServe()
}

func (i *internalHttpServer) Stop() error {
	return i.server.Shutdown(context.Background())
}

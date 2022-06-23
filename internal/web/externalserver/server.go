package externalserver

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	s "tec-doc/internal/service"
)

type Server interface {
	Start() error
	Stop() error
}

type externalHttpServer struct {
	router  *gin.Engine
	server  http.Server
	service *s.Service
}

func New(bindingAddress string, service *s.Service) *externalHttpServer {
	router := gin.Default()
	serv := &externalHttpServer{
		router:  router,
		service: service,
		server: http.Server{
			Addr:    bindingAddress,
			Handler: router,
		},
	}
	serv.configureRouter()
	return serv
}

func (e *externalHttpServer) configureRouter() {
	e.router.GET("/excel_template", e.ExcelTemplate)
	e.router.POST("/load_from_excel", e.LoadFromExcel)
}

func (e *externalHttpServer) Start() error {
	log.Info().Msg("start external server on " + e.server.Addr)
	return e.server.ListenAndServe()
}

func (e *externalHttpServer) Stop() error {
	return e.server.Shutdown(context.Background())
}

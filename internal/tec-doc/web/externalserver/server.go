package externalserver

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
	"tec-doc/internal/tec-doc/model"
	"tec-doc/pkg/metrics"
)

type Service interface {
	ExcelTemplateForClient() ([]byte, error)
	AddFromExcel(ctx *gin.Context, products []model.Product, supplierID int64, userID int64) error
	GetSupplierTaskHistory(ctx context.Context, supplierID int64, limit int, offset int) ([]model.Task, error)
	GetProductsHistory(ctx context.Context, uploadID int64, limit int, offset int) ([]model.Product, error)
	GetArticles(ctx context.Context, dataSupplierID int, article string) ([]model.Article, error)
	GetBrand(ctx context.Context, brandName string) (*model.Brand, error)
}

type Server interface {
	Start() error
	Stop() error
}

type externalHttpServer struct {
	router  *gin.Engine
	server  http.Server
	metrics *metrics.Metrics
	service Service
	logger  *zerolog.Logger
}

func (e *externalHttpServer) Start() error {
	return e.server.ListenAndServe()
}

func (e *externalHttpServer) Stop() error {
	return e.server.Shutdown(context.Background())
}

func (e *externalHttpServer) configureRouter() {
	e.router.Use(gin.Recovery())
	e.router.Use(e.Authorize)
	e.router.Use(e.MiddleWareMetric)
	e.router.GET("/excel_template", e.ExcelTemplate)
	e.router.POST("/load_from_excel", e.LoadFromExcel)
	e.router.GET("/task_history", e.GetSupplierTaskHistory)
	e.router.GET("/product_history", e.GetProductsHistory)
	e.router.GET("/tecdoc_articles", e.GetTecDocArticles)
}

func New(bindingPort string, service Service, logger *zerolog.Logger, mts *metrics.Metrics) *externalHttpServer {
	router := gin.Default()
	serv := &externalHttpServer{
		router:  router,
		service: service,
		logger:  logger,
		metrics: mts,
		server: http.Server{
			Addr:    bindingPort,
			Handler: router,
		},
	}
	serv.configureRouter()
	return serv
}

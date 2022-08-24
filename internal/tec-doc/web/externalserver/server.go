package externalserver

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
	_ "tec-doc/docs"
	"tec-doc/internal/tec-doc/config"
	"tec-doc/internal/tec-doc/model"
	"tec-doc/pkg/clients/services"
	"tec-doc/pkg/metrics"
)

type Service interface {
	ExcelTemplateForClient() ([]byte, error)
	AddFromExcel(ctx *gin.Context, products []model.Product, supplierID int64, userID int64) error
	GetSupplierTaskHistory(ctx context.Context, supplierID int64, limit int, offset int) ([]model.Task, error)
	GetProductsHistory(ctx context.Context, uploadID int64, limit int, offset int) ([]model.Product, error)
	GetArticles(dataSupplierID int, article string) ([]model.Article, error)
	GetBrand(brandName string) (*model.Brand, error)

	Scope() *config.Scope
	Abac() services.ABAC
	Suppliers() services.Suppliers
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
	e.router.Use(e.MiddleWareMetric)
	api := e.router.Group("/api")
	{
		//api.Use(e.Authorize)
		api.GET("/excel_template", e.ExcelTemplate)
		api.POST("/load_from_excel", e.LoadFromExcel)
		api.GET("/task_history", e.GetSupplierTaskHistory)
		api.POST("/product_history", e.GetProductsHistory)
		api.GET("/tecdoc_articles", e.GetTecDocArticles)
	}
}

func New(bindingPort string, service Service, logger *zerolog.Logger, mts *metrics.Metrics) *externalHttpServer {
	router := gin.Default()
	serv := &externalHttpServer{
		router:  router,
		service: service,
		logger:  logger,
		metrics: mts,
		server: http.Server{
			Addr: bindingPort,
		},
	}
	serv.configureRouter()
	return serv
}

package externalserver

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	_ "tec-doc/docs"
	"tec-doc/internal/tec-doc/config"
	"tec-doc/pkg/clients/services"
	"tec-doc/pkg/metrics"
	"tec-doc/pkg/model"
)

type Service interface {
	GetProductsEnrichedExcel(products []model.Product) (data []byte, err error)
	ExcelTemplateForClient() ([]byte, error)
	AddFromExcel(ctx *gin.Context, products []model.Product, supplierID int64, userID int64) error
	LoadFromExcel(bodyData io.Reader) (products []model.Product, err error)
	GetSupplierTaskHistory(ctx context.Context, supplierID int64, limit int, offset int) ([]model.Task, error)
	GetProductsHistory(ctx context.Context, uploadID int64, limit int, offset int) ([]model.Product, error)
	ExcelProductsHistoryWithStatus(ctx context.Context, uploadID, status int64) ([]byte, error)
	GetArticles(dataSupplierID int, article string) ([]model.Article, error)
	GetBrand(brandName string) (*model.Brand, error)
	Enrichment(product []model.Product) (productsEnriched []model.ProductEnriched, err error)

	Scope() *config.Scope
	Abac() services.ABAC
	Suppliers() services.Suppliers
}

type Server interface {
	Start() error
	Stop() error
}

type externalHttpServer struct {
	testMode bool
	router   *gin.Engine
	server   http.Server
	metrics  *metrics.Metrics
	service  Service
	logger   *zerolog.Logger
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
	api := e.router.Group("/api/v1")
	{
		api.Use(e.Authorize)
		api.GET("/excel", e.ExcelTemplate)
		api.POST("/excel", e.LoadFromExcel)
		api.POST("/excel/products/enrichment", e.GetProductsEnrichedExcel)
		api.POST("/excel/products/errors", e.ExcelProductsWithErrors)
		api.GET("/task_history", e.GetSupplierTaskHistory)
		api.POST("/product_history", e.GetProductsHistory)
		api.POST("/articles/enrichment", e.GetTecDocArticles)
	}
}

func New(bindingPort string, service Service, logger *zerolog.Logger, mts *metrics.Metrics, testMode bool) Server {
	router := gin.Default()
	serv := &externalHttpServer{
		testMode: testMode,
		router:   router,
		service:  service,
		logger:   logger,
		metrics:  mts,
		server: http.Server{
			Addr:    bindingPort,
			Handler: router,
		},
	}
	serv.configureRouter()
	return serv
}

package externalserver

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"os"
	_ "tec-doc/docs"
	"tec-doc/internal/tec-doc/config"
	"tec-doc/pkg/clients/services"
	"tec-doc/pkg/ginLogger"
	"tec-doc/pkg/metrics"
	"tec-doc/pkg/model"
)

type Service interface {
	GetProductsEnrichedExcel(products []model.Product) ([]byte, error)
	ExcelTemplateForClient() ([]byte, error)
	AddFromExcel(ctx *gin.Context, products []model.Product, supplierID int64, userID int64) error
	LoadFromExcel(bodyData io.Reader) (products []model.Product, err error)
	GetSupplierTaskHistory(ctx context.Context, supplierID int64, limit int, offset int) ([]model.Task, error)
	GetProductsHistory(ctx context.Context, uploadID int64, limit int, offset int) ([]model.Product, error)
	ExcelProductsHistoryWithStatus(ctx context.Context, uploadID, status int64) ([]byte, error)
	GetArticles(dataSupplierID int, article string) ([]model.Article, error)
	GetBrand(brandName string) (*model.Brand, error)
	Enrichment(product []model.Product) (productsEnriched []model.ProductEnriched)

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
	e.router.Use(ginLogger.Logger(os.Stdout))
	api := e.router.Group("/api/v1")
	{
		api.Use(e.Authorize)
		api.GET("/excel", e.ExcelTemplate)
		api.POST("/excel", e.LoadFromExcel)
		api.POST("/excel/products/errors", e.ExcelProductsWithErrors)
		api.GET("/history/task", e.GetSupplierTaskHistory)
		api.POST("/history/product", e.GetProductsHistory)
		api.POST("/articles/enrichment", e.GetTecDocArticles)
	}

	admin := e.router.Group("/api/v1/admin")
	admin.POST("/excel/products/enrichment", e.GetProductsEnrichedExcel)
}

func New(bindingPort string, service Service, logger *zerolog.Logger, mts *metrics.Metrics, testMode bool) Server {
	router := gin.New()
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

package externalserver

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
	"tec-doc/internal/model"
	"tec-doc/internal/store/postgres"
	"tec-doc/internal/web/externalserver/middleware"
	m "tec-doc/internal/web/metrics"
	"time"
)

type Service interface {
	ExcelTemplateForClient() ([]byte, error)
	AddFromExcel(bodyData io.Reader, ctx *gin.Context) error
	GetSupplierTaskHistory(ctx context.Context, supplierID int64, limit int, offset int) ([]model.Task, error)
	GetProductsHistory(ctx context.Context, tx postgres.Transaction, uploadID int64, limit int, offset int) ([]model.Product, error)
	//...
}

type Server interface {
	Start() error
	Stop() error
}

type externalHttpServer struct {
	router  *gin.Engine
	server  http.Server
	metrics *m.Metrics
	service Service
	logger  *zerolog.Logger
}

func New(bindingAddress string, service Service, logger *zerolog.Logger) *externalHttpServer {
	router := gin.Default()
	serv := &externalHttpServer{
		router:  router,
		service: service,
		logger:  logger,
		metrics: m.NewMetrics("external", "HttpServer"),
		server: http.Server{
			Addr:    bindingAddress,
			Handler: router,
		},
	}
	serv.configureRouter()
	return serv
}

func (e *externalHttpServer) configureRouter() {
	e.router.Use(middleware.Authorize)
	e.router.Use(gin.Recovery())
	e.router.Use(e.MiddleWareMetric)
	e.router.GET("/excel_template", e.ExcelTemplate)
	e.router.POST("/load_from_excel", e.LoadFromExcel)
	e.router.GET("/task_history", e.GetSupplierTaskHistory)
	e.router.GET("/product_history", e.ProductHistory)
}

func (e *externalHttpServer) Start() error {
	log.Info().Msg("start external server on " + e.server.Addr)
	return e.server.ListenAndServe()
}

func (e *externalHttpServer) Stop() error {
	return e.server.Shutdown(context.Background())
}

func (e *externalHttpServer) MiddleWareMetric(c *gin.Context) {
	t := time.Now()
	c.Next()
	status := strconv.Itoa(c.Writer.Status())
	e.metrics.Collector.WithLabelValues(
		m.InternalServerComponent,
		c.Request.Method,
		c.Request.URL.Path,
		status,
	).Inc()

	defer func() {
		e.metrics.LeadTime.WithLabelValues(
			m.InternalServerComponent,
			c.Request.Method,
			c.Request.URL.Path,
			strconv.FormatInt(time.Since(t).Milliseconds(), 10),
		).Observe(float64(time.Since(t).Milliseconds()))
	}()

	defer func() {
		e.metrics.LeadTimeQua.WithLabelValues(
			m.InternalServerComponent,
			c.Request.Method,
			c.Request.URL.Path,
			strconv.FormatInt(time.Since(t).Milliseconds(), 10),
		).Observe(float64(time.Since(t).Milliseconds()))
	}()

	e.metrics.Rating.WithLabelValues(
		m.InternalServerComponent,
		c.Request.Method,
		c.Request.URL.Path,
		status,
	).Inc()
}

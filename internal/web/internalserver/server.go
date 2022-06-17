package internalserver

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tec-doc/internal/web"
	m "tec-doc/internal/web/metrics"
	"time"
)

type internalHttpServer struct {
	router  *gin.Engine
	server  *http.Server
	metrics *m.Metrics
}

func NewInternalServer(bindingAddress string) web.Server {
	serv := new(internalHttpServer)
	serv.metrics = m.NewMetrics("internal", "HttpServer")

	serv.router = gin.Default()

	serv.router.Use(gin.Recovery())
	serv.router.Use(serv.MiddleWareMetric)

	serv.router.GET("/helth", serv.Helth)
	serv.router.GET("/readiness", serv.Readiness)
	serv.router.GET("/metrics", serv.Metrics)

	serv.server = &http.Server{
		Addr:    bindingAddress,
		Handler: serv.router,
	}
	return serv
}

func (i *internalHttpServer) Start() error {
	return i.server.ListenAndServe()
}

func (i *internalHttpServer) Stop() error {
	return i.server.Shutdown(context.Background())
}

func (i *internalHttpServer) MiddleWareMetric(c *gin.Context) {
	t := time.Now()
	c.Next()
	status := strconv.Itoa(c.Writer.Status())
	i.metrics.Collector.WithLabelValues(
		m.InternalServerComponent,
		c.Request.Method,
		c.Request.URL.Path,
		status,
	).Inc()

	defer func() {
		i.metrics.LeadTime.WithLabelValues(
			m.InternalServerComponent,
			c.Request.Method,
			c.Request.URL.Path,
			strconv.FormatInt(time.Since(t).Milliseconds(), 10),
		).Observe(float64(time.Since(t).Milliseconds()))
	}()

	defer func() {
		i.metrics.LeadTimeQua.WithLabelValues(
			m.InternalServerComponent,
			c.Request.Method,
			c.Request.URL.Path,
			strconv.FormatInt(time.Since(t).Milliseconds(), 10),
		).Observe(float64(time.Since(t).Milliseconds()))
	}()

	i.metrics.Rating.WithLabelValues(
		m.InternalServerComponent,
		c.Request.Method,
		c.Request.URL.Path,
		status,
	).Inc()
}

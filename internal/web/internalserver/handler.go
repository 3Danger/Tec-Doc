package internalserver

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (i *internalHttpServer) Health(c *gin.Context) {
	c.Status(200)
}

func (i *internalHttpServer) Readiness(c *gin.Context) {
	c.Status(200)
}

func (i *internalHttpServer) Metrics(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

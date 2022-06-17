package internalserver

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (i *internalHttpServer) Helth(c *gin.Context) {
	c.Writer.WriteHeader(200)
	c.Writer.WriteHeaderNow()

}

func (i *internalHttpServer) Readiness(c *gin.Context) {
	c.Writer.WriteHeader(200)
	c.Writer.WriteHeaderNow()
}

func (i *internalHttpServer) Metrics(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

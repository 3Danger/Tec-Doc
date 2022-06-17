package internalserver

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func initInternalRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/helth", helth)
	router.GET("/readiness", readiness)
	router.GET("/metrics", metrics)
	return router
}

func helth(c *gin.Context) {
	c.Writer.WriteHeader(200)
	c.Writer.WriteHeaderNow()

}

func readiness(c *gin.Context) {
	c.Writer.WriteHeader(200)
	c.Writer.WriteHeaderNow()
}

func metrics(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

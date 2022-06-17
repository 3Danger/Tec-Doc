package internalserver

import "github.com/gin-gonic/gin"

///helth, /readiness, /metrics
func initRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/helth", helth)
	router.GET("/readiness", readiness)
	router.GET("/metrics", metrics)
	return router
}

func helth(c *gin.Context) {
	//c.String(200, "HELTH")
	c.Writer.WriteHeader(200)
	c.Writer.WriteHeaderNow()

}

func readiness(c *gin.Context) {
	//c.String(200, "READINESS")
	c.Writer.WriteHeader(200)
	c.Writer.WriteHeaderNow()
}

func metrics(c *gin.Context) {
	c.String(200, "METRICS")
}

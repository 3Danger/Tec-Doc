package internalserver

import "github.com/gin-gonic/gin"

const ERROR = "error"
const MESSAGE = "message"

func jsonError(err error, status int, c *gin.Context) {
	c.JSON(status, gin.H{
		ERROR: err.Error(),
	})
}

func jsonMessage(msg string, status int, c *gin.Context) {
	c.JSON(status, gin.H{
		MESSAGE: msg,
	})
}

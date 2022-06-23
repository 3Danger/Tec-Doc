package utils

import "github.com/gin-gonic/gin"

const ERROR = "error"
const MESSAGE = "message"

func JsonError(err error, status int, c *gin.Context) {
	c.JSON(status, gin.H{
		ERROR: err.Error(),
	})
}

func JsonMessage(msg string, status int, c *gin.Context) {
	c.JSON(status, gin.H{
		MESSAGE: msg,
	})
}

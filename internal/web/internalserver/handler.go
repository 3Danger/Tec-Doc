package internalserver

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
)

func (i *internalHttpServer) Helth(c *gin.Context) {
	defer func() { _ = c.Request.Body.Close() }()
	c.Status(200)
}

func (i *internalHttpServer) Readiness(c *gin.Context) {
	defer func() { _ = c.Request.Body.Close() }()
	c.Status(200)
}

func (i *internalHttpServer) Metrics(c *gin.Context) {
	defer func() { _ = c.Request.Body.Close() }()
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

func (i *internalHttpServer) ExcelTemplate(c *gin.Context) {
	defer func() { _ = c.Request.Body.Close() }()
	excelTemplate, err := i.service.ExcelTemplateForClient("")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.Data(200, "application/vnd.ms-excel", excelTemplate)
	}
}

func (i *internalHttpServer) LoadFromExcel(c *gin.Context) {
	defer func() { _ = c.Request.Body.Close() }()
	rawXls, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	products, err := i.service.LoadFromExcel(rawXls)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"msg": "success"})
	}
	_ = products
}

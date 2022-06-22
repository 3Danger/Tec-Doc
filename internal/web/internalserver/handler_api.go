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
	file, header, err := c.Request.FormFile("excel_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	defer func() { _ = file.Close() }()
	_ = header
	rawXls, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
		return
	}
	products, err := i.service.LoadFromExcel(rawXls)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Error": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "success"})
	}
	_ = products
}

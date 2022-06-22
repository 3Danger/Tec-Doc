package internalserver

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
)

const ContentTypeExcel = "application/vnd.ms-excel"

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
		jsonError(err, http.StatusInternalServerError, c)
	} else {
		c.Data(200, ContentTypeExcel, excelTemplate)
	}
}

func (i *internalHttpServer) LoadFromExcel(c *gin.Context) {
	defer func() { _ = c.Request.Body.Close() }()
	rawXls, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		jsonError(err, http.StatusBadRequest, c)
	} else {
		products, err := i.service.LoadFromExcel(rawXls)
		if err != nil {
			jsonError(err, http.StatusInternalServerError, c)
		} else {
			jsonMessage("success", http.StatusOK, c)
		}
		_ = products
	}
}

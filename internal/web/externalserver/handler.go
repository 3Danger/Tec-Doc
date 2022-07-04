package externalserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tec-doc/internal/web/externalserver/middleware"
)

const ContentTypeExcel = "application/vnd.ms-excel"

func (e *externalHttpServer) ExcelTemplate(c *gin.Context) {
	excelTemplate, err := e.service.ExcelTemplateForClient()
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.Data(200, ContentTypeExcel, excelTemplate)
}

func (e *externalHttpServer) LoadFromExcel(c *gin.Context) {
	err := e.service.AddFromExcel(c.Request.Body, c)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

//func (e *externalHttpServer) ProductHistory(c *gin.Context) {
//	var t int64 = 1
//	c.Set("upload_id", t)
//	c.Set("limit", 10)
//	c.Set("offset", 0)
//	productHistory, err := e.service.GetProductHistory(c)
//	if err != nil {
//		e.logger.Error().Err(err).Send()
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}
//	c.JSON(200, productHistory)
//}

func (e *externalHttpServer) GetProductsHistory(c *gin.Context) {
	uploadID, err := strconv.ParseInt(c.Request.Header.Get("upload_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "can't get limit",
		})
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "can't get limit",
		})
		return
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "can't get offset",
		})
		return
	}

	productsHistory, err := e.service.GetProductsHistory(c, nil, uploadID, limit, offset)
	c.JSON(http.StatusOK, productsHistory)
}

func (e *externalHttpServer) GetSupplierTaskHistory(c *gin.Context) {
	supplierID, _, err := middleware.CredentialsFromContext(c)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "can't get limit",
		})
		return
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "can't get offset",
		})
		return
	}

	rawTasks, err := e.service.GetSupplierTaskHistory(c, supplierID, limit, offset)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, rawTasks)
}

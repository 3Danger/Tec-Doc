package externalserver

import (
	"encoding/json"
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

func (e *externalHttpServer) GetProductsHistory(c *gin.Context) {

	type ReqStruct struct {
		UploadID int64 `json:"upload_id"`
	}

	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()

	var rs ReqStruct

	if err := dec.Decode(&rs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "can't get upload_id",
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

	productsHistory, err := e.service.GetProductsHistory(c, nil, rs.UploadID, limit, offset)
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

	rawTasks, err := e.service.GetSupplierTaskHistory(c, nil, supplierID, limit, offset)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, rawTasks)
}

func (e *externalHttpServer) GetTecDocArticles(c *gin.Context) {
	type ReqStruct struct {
		ArticleNumber string `json:"ArticleNumber"`
		Brand         string `json:"Brand"`
	}

	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()

	var rs ReqStruct

	if err := dec.Decode(&rs); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "can't get brand and article number",
		})
		return
	}

	brand := e.serviceю
}

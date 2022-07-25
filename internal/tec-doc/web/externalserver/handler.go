package externalserver

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tec-doc/pkg/errinfo"
)

func (e *externalHttpServer) ExcelTemplate(c *gin.Context) {
	excelTemplate, err := e.service.ExcelTemplateForClient()
	if err != nil {
		msg, status := errinfo.GetErrorInfo(errinfo.InvalidExcelData)
		e.logger.Error().Err(err).Send()
		c.JSON(status, msg)
		return
	}
	c.Data(http.StatusOK, "application/vnd.ms-excel", excelTemplate)
}

func (e *externalHttpServer) LoadFromExcel(c *gin.Context) {
	supplierID, userID := c.GetInt64("X-Supplier-Id"), c.GetInt64("X-User-Id")

	products, err := e.loadFromExcel(c.Request.Body)
	if err != nil {
		msg, status := errinfo.GetErrorInfo(errinfo.InvalidSupplierID)
		e.logger.Error().Err(err).Send()
		c.JSON(status, msg)
		return
	}

	err = e.service.AddFromExcel(c, products, supplierID, userID)
	if err != nil {
		msg, status := errinfo.GetErrorInfo(errinfo.InternalServerErr)
		e.logger.Error().Err(err).Send()
		c.JSON(status, msg)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "данные успешно загружены",
	})
}

func (e *externalHttpServer) GetProductsHistory(c *gin.Context) {
	var rs map[string]int64
	if err := json.NewDecoder(c.Request.Body).Decode(&rs); err != nil {
		msg, status := errinfo.GetErrorInfo(errinfo.InvalidTaskID)
		e.logger.Error().Err(err).Send()
		c.JSON(status, msg)
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		msg, status := errinfo.GetErrorInfo(errinfo.InvalidLimit)
		e.logger.Error().Err(err).Send()
		c.JSON(status, msg)
		return
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		msg, status := errinfo.GetErrorInfo(errinfo.InvalidOffset)
		e.logger.Error().Err(err).Send()
		c.JSON(status, msg)
		return
	}

	productsHistory, err := e.service.GetProductsHistory(c, nil, rs["upload_id"], limit, offset)
	if err != nil {
		msg, status := errinfo.GetErrorInfo(errinfo.InternalServerErr)
		e.logger.Error().Err(err).Send()
		c.JSON(status, msg)
		return
	}
	c.JSON(http.StatusOK, productsHistory)
}

func (e *externalHttpServer) GetSupplierTaskHistory(c *gin.Context) {
	supplierID, _, err := CredentialsFromContext(c)
	if err != nil {
		msg, status := errinfo.GetErrorInfo(errinfo.InvalidTaskID)
		e.logger.Error().Err(err).Send()
		c.JSON(status, msg)
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		msg, status := errinfo.GetErrorInfo(errinfo.InvalidLimit)
		e.logger.Error().Err(err).Send()
		c.JSON(status, msg)
		return
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		msg, status := errinfo.GetErrorInfo(errinfo.InvalidOffset)
		e.logger.Error().Err(err).Send()
		c.JSON(status, msg)
		return
	}

	rawTasks, err := e.service.GetSupplierTaskHistory(c, nil, supplierID, limit, offset)
	if err != nil {
		msg, status := errinfo.GetErrorInfo(errinfo.InternalServerErr)
		e.logger.Error().Err(err).Send()
		c.JSON(status, msg)
		return
	}

	c.JSON(http.StatusOK, rawTasks)
}

func (e *externalHttpServer) GetTecDocArticles(c *gin.Context) {
	var rs map[string]string
	if err := json.NewDecoder(c.Request.Body).Decode(&rs); err != nil {
		msg, status := errinfo.GetErrorInfo(errinfo.InvalidTecDocParams)
		e.logger.Error().Err(err).Send()
		c.JSON(status, msg)
		return
	}

	brand, err := e.service.GetBrand(c, rs["Brand"])
	if err != nil {
		msg, status := errinfo.GetErrorInfo(errinfo.InternalServerErr)
		e.logger.Error().Err(err).Send()
		c.JSON(status, msg)
		return
	}

	articles, err := e.service.GetArticles(c, brand.SupplierId, rs["ArticleNumber"])
	if err != nil {
		msg, status := errinfo.GetErrorInfo(errinfo.InternalServerErr)
		e.logger.Error().Err(err).Send()
		c.JSON(status, msg)
		return
	}

	c.JSON(http.StatusOK, articles)
}

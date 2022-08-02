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
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidExcelData))
		return
	}
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelTemplate)
}

func (e *externalHttpServer) LoadFromExcel(c *gin.Context) {
	supplierID, userID := c.GetInt64("X-Supplier-Id"), c.GetInt64("X-User-Id")

	products, err := e.loadFromExcel(c.Request.Body)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidSupplierID))
		return
	}

	err = e.service.AddFromExcel(c, products, supplierID, userID)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InternalServerErr))
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "данные успешно загружены",
	})
}

func (e *externalHttpServer) GetProductsHistory(c *gin.Context) {
	var rs map[string]int64
	if err := json.NewDecoder(c.Request.Body).Decode(&rs); err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidTaskID))
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidLimit))
		return
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidOffset))
		return
	}

	productsHistory, err := e.service.GetProductsHistory(c, rs["upload_id"], limit, offset)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InternalServerErr))
		return
	}
	c.JSON(http.StatusOK, productsHistory)
}

func (e *externalHttpServer) GetSupplierTaskHistory(c *gin.Context) {
	supplierID, _, err := CredentialsFromContext(c)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidTaskID))
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidLimit))
		return
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidOffset))
		return
	}

	rawTasks, err := e.service.GetSupplierTaskHistory(c, supplierID, limit, offset)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InternalServerErr))
		return
	}

	c.JSON(http.StatusOK, rawTasks)
}

func (e *externalHttpServer) GetTecDocArticles(c *gin.Context) {
	var rs map[string]string
	if err := json.NewDecoder(c.Request.Body).Decode(&rs); err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidTecDocParams))
		return
	}

	brand, err := e.service.GetBrand(c, rs["Brand"])
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InternalServerErr))
		return
	}

	articles, err := e.service.GetArticles(c, brand.SupplierId, rs["ArticleNumber"])
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InternalServerErr))
		return
	}

	c.JSON(http.StatusOK, articles)
}

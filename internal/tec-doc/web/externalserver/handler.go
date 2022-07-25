package externalserver

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// todo: сделать аналоги ошибок на русском
func (e *externalHttpServer) ExcelTemplate(c *gin.Context) {
	excelTemplate, err := e.service.ExcelTemplateForClient()
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "внутренняя ошибка сервера",
		})
		return
	}
	c.Data(http.StatusOK, "application/vnd.ms-excel", excelTemplate)
}

// todo: получение header, query, params url, body делаем на уровне web server
func (e *externalHttpServer) LoadFromExcel(c *gin.Context) {
	supplierID, userID := c.GetInt64("X-Supplier-Id"), c.GetInt64("X-User-Id")

	products, err := e.loadFromExcel(c.Request.Body)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректные данные в excel",
		})
		return
	}

	err = e.service.AddFromExcel(c, products, supplierID, userID)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "внутренняя ошибка сервера",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "данные успешно загружены",
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
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректный ID задания",
		})
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректный параметр limit",
		})
		return
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректный параметр offset",
		})
		return
	}

	productsHistory, err := e.service.GetProductsHistory(c, nil, rs.UploadID, limit, offset)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "внутренняя ошибка сервера",
		})
		return
	}
	c.JSON(http.StatusOK, productsHistory)
}

func (e *externalHttpServer) GetSupplierTaskHistory(c *gin.Context) {
	supplierID, _, err := CredentialsFromContext(c)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "некорректный ID задания",
		})
		return
	}

	limit, err := strconv.Atoi(c.Query("limit"))
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректный параметр limit",
		})
		return
	}

	offset, err := strconv.Atoi(c.Query("offset"))
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректный параметр offset",
		})
		return
	}

	rawTasks, err := e.service.GetSupplierTaskHistory(c, nil, supplierID, limit, offset)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "внутренняя ошибка сервера",
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
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "некорректное имя бренда или номер артикула",
		})
		return
	}

	brand, err := e.service.GetBrand(c, rs.Brand)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "внутренняя ошибка сервера",
		})
		return
	}

	articles, err := e.service.GetArticles(c, brand.SupplierId, rs.ArticleNumber)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "внутренняя ошибка сервера",
		})
		return
	}

	c.JSON(http.StatusOK, articles)
}

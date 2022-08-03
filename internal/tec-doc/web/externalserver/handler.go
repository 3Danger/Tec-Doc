package externalserver

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strconv"
	"tec-doc/pkg/errinfo"
)

// @Summary ExcelTemplate
// @Tags excel
// @Description download excel table template
// @ID excel_template
// @Produce application/vnd.ms-excel
// @Param X-User-Id header string true "ID of user"
// @Param X-Supplier-Id header string true "ID of supplier"
// @Success 200 {array} byte
// @Failure 500 {object} errinfo.errInf
// @Router /excel_template [get]
func (e *externalHttpServer) ExcelTemplate(c *gin.Context) {
	excelTemplate, err := e.service.ExcelTemplateForClient()
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidExcelData))
		return
	}
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excelTemplate)
}

// @Summary LoadFromExcel
// @Tags excel
// @Description for upload excel table with products into
// @ID load_from_excel
// @Produce json
// @Param excel_file formData file true "excel file"
// @Param X-User-Id header string true "ID of user"
// @Param X-Supplier-Id header string true "ID of supplier"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /load_from_excel [post]
func (e *externalHttpServer) LoadFromExcel(c *gin.Context) {
	supplierID, userID := c.GetInt64("X-Supplier-Id"), c.GetInt64("X-User-Id")

	file, _, err := c.Request.FormFile("excel_file")
	if err != nil {
		log.Error().Err(errinfo.InvalidNotFile).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidNotFile))
		return
	}
	defer func() { _ = file.Close() }()
	products, err := e.loadFromExcel(file)
	if err != nil {
		e.logger.Error().Err(err).Send()
		if err.Error() == "empty data" || err == io.EOF {
			c.JSON(errinfo.GetErrorInfo(errinfo.InvalidExcelEmpty))
			return
		}
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidExcelData))
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

type RequestBody struct {
	UploadID int64 `json:"uploadID"`
}

// @Summary GetProductsHistory
// @Tags product
// @Description getting product list
// @ID products_history
// @Accept json
// @Produce json
// @Param limit query string true "limit of contents"
// @Param offset query string true "offset of contents"
// @Param X-User-Id header string true "ID of user"
// @Param X-Supplier-Id header string true "ID of supplier"
// @Param RequestBody body RequestBody true "The input body struct"
// @Success 200 {array} model.Product
// @Failure 500 {object} errinfo.errInf
// @Router /product_history [get]
func (e *externalHttpServer) GetProductsHistory(c *gin.Context) {
	type RequestBody struct {
		UploadID int64 `json:"uploadID"`
	}
	var rq RequestBody

	if err := json.NewDecoder(c.Request.Body).Decode(&rq); err != nil {
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

	productsHistory, err := e.service.GetProductsHistory(c, rq.UploadID, limit, offset)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InternalServerErr))
		return
	}
	c.JSON(http.StatusOK, productsHistory)
}

// @Summary GetSupplierTaskHistory
// @Tags product
// @Description getting task list
// @ID supplier_task_history
// @Produce json
// @Param limit query string true "limit of contents"
// @Param offset query string true "offset of contents"
// @Param X-User-Id header string true "ID of user"
// @Param X-Supplier-Id header string true "ID of supplier"
// @Success 200 {array} model.Task
// @Failure 500 {object} errinfo.errInf
// @Router /task_history [get]
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
	type Request struct {
		Brand         string `json:"Brand"`
		ArticleNumber string `json:"ArticleNumber"`
	}
	var rq Request

	if err := json.NewDecoder(c.Request.Body).Decode(&rq); err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidTecDocParams))
		return
	}

	brand, err := e.service.GetBrand(c, rq.Brand)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InternalServerErr))
		return
	}

	articles, err := e.service.GetArticles(c, brand.SupplierId, rq.ArticleNumber)
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InternalServerErr))
		return
	}

	c.JSON(http.StatusOK, articles)
}

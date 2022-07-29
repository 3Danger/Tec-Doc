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
// @Param X-User-Id query string true "ID of user"
// @Param X-Supplier-Id query string true "ID of supplier"
// @Success 200 {array} byte
// @Failure 500 {object} errinfo.errInf
// @Router /excel_template [get]
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

// @Summary LoadFromExcel
// @Tags excel
// @Description for upload excel table with products into
// @ID load_from_excel
// @Accept application/vnd.ms-excel
// @Produce json
// @Param X-User-Id query string true "ID of user"
// @Param X-Supplier-Id query string true "ID of supplier"
// @Param file formData file true "excel file"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /load_from_excel [post]
func (e *externalHttpServer) LoadFromExcel(c *gin.Context) {
	supplierID, userID := c.GetInt64("X-Supplier-Id"), c.GetInt64("X-User-Id")

	if c.Request.Body == http.NoBody {
		msg, status := errinfo.GetErrorInfo(errinfo.InvalidBodyEmpty)
		log.Error().Err(errinfo.InvalidBodyEmpty).Send()
		c.JSON(status, msg)
		return
	}
	products, err := e.loadFromExcel(c.Request.Body)
	if err != nil {
		msg, status := errinfo.GetErrorInfo(errinfo.InvalidExcelData)
		if err.Error() == "empty data" || err == io.EOF {
			msg, status = errinfo.GetErrorInfo(errinfo.InvalidExcelEmpty)
		}
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
	c.JSON(http.StatusOK, "данные успешно загружены")
}

// @Summary GetProductsHistory
// @Tags product
// @Description getting product list
// @ID products_history
// @Accept json
// @Produce json
// @Param upload_id body object true "ID of the task sender"
// @Param limit query string true "limit of contents"
// @Param offset query string true "offset of contents"
// @Param X-User-Id query string true "ID of user"
// @Param X-Supplier-Id query string true "ID of supplier"
// @Success 200 {array} model.Product
// @Failure 500 {object} errinfo.errInf
// @Router /product_history [get]
func (e *externalHttpServer) GetProductsHistory(c *gin.Context) {
	rs := struct {
		UploadId int64 `json:"upload_id" required:"true"`
	}{}

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

	productsHistory, err := e.service.GetProductsHistory(c, nil, rs.UploadId, limit, offset)
	if err != nil {
		msg, status := errinfo.GetErrorInfo(errinfo.InternalServerErr)
		e.logger.Error().Err(err).Send()
		c.JSON(status, msg)
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
// @Param X-User-Id query string true "ID of user"
// @Param X-Supplier-Id query string true "ID of supplier"
// @Success 200 {array} model.Task
// @Failure 500 {object} errinfo.errInf
// @Router /task_history [get]
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

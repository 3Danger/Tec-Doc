package externalserver

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
	"tec-doc/pkg/errinfo"
	"tec-doc/pkg/model"
)

// @Summary ExcelTemplate
// @Tags excel
// @Description download excel table template
// @ID excel_template
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
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

// @Summary ProductsEnrichedExcel
// @Tags excel
// @Description Enrichment excel file, limit entiies in file = 10000
// @Produce json
// @ID enrich_excel
// @Param excel_file body []byte true "binary excel file"
// @Success 200 {array} byte
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /excel/enrichment [post]
func (e *externalHttpServer) GetProductsEnrichedExcel(c *gin.Context) {
	var (
		err      error
		products []model.Product
	)

	if products, err = e.loadFromExcel(c.Request.Body); err != nil {
		e.logger.Error().Err(err).Send()
		if err.Error() == "empty data" || err == io.EOF {
			c.JSON(errinfo.GetErrorInfo(errinfo.InvalidExcelEmpty))
			return
		}
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidExcelData))
		return
	}

	if len(products) > 10000 {
		e.logger.Error().Err(errinfo.InvalidExcelLimit).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidExcelLimit))
		return
	}

	excel, err := e.service.GetProductsEnrichedExcel(products)
	if err != nil {
		e.logger.Error().Err(err).Msg("can't to create excel enrichment file")
		c.JSON(errinfo.GetErrorInfo(errinfo.InternalServerErr))
		return
	}
	//{ //Посмотреть содержимое файла без танцев с бубном
	//	dir, _ := os.UserHomeDir()
	//	_ = ioutil.WriteFile(dir + "/excel_test_file.xlsx", excel, 0644)
	//}
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", excel)
}

// @Summary LoadFromExcel
// @Tags excel
// @Description upload excel table containing products info
// @ID load_from_excel
// @Param excel_file body []byte true "binary excel file"
// @Param X-User-Id header string true "ID of user"
// @Param X-Supplier-Id header string true "ID of supplier"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /load_from_excel [post]
func (e *externalHttpServer) LoadFromExcel(c *gin.Context) {
	supplierID, userID := c.GetInt64("X-Supplier-Id"), c.GetInt64("X-User-Id")
	products, err := e.loadFromExcel(c.Request.Body)
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

// @Summary GetProductsHistory
// @Tags product
// @Description get product list
// @ID products_history
// @Accept json
// @Produce json
// @Param limit query string true "limit of contents"
// @Param offset query string true "offset of contents"
// @Param X-User-Id header string true "ID of user"
// @Param X-Supplier-Id header string true "ID of supplier"
// @Param InputBody body model.UploadIdRequest true "The input body.<br /> UploadID is ID of previously uploaded task."
// @Success 200 {array} model.Product
// @Failure 500 {object} errinfo.errInf
// @Router /product_history [post]
func (e *externalHttpServer) GetProductsHistory(c *gin.Context) {
	var rq model.UploadIdRequest

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
// @Description get task list
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
	var rq model.GetTecDocArticlesRequest

	if err := json.NewDecoder(c.Request.Body).Decode(&rq); err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidTecDocParams))
		return
	}

	brand, err := e.service.GetBrand(rq.Brand)
	if err != nil {
		e.logger.Error().Err(err).Send()
		if errors.Is(err, errinfo.NoTecDocBrandFound) {
			c.JSON(errinfo.GetErrorInfo(err))
		} else {
			c.JSON(errinfo.GetErrorInfo(errinfo.InternalServerErr))
		}
		return
	}

	articles, err := e.service.GetArticles(brand.SupplierId, rq.ArticleNumber)
	if err != nil {
		e.logger.Error().Err(err).Send()
		if errors.Is(err, errinfo.NoTecDocArticlesFound) {
			c.JSON(errinfo.GetErrorInfo(err))
		} else {
			c.JSON(errinfo.GetErrorInfo(errinfo.InternalServerErr))
		}
		return
	}

	c.JSON(http.StatusOK, articles)
}

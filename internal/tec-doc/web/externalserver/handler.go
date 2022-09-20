package externalserver

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	exl "github.com/xuri/excelize/v2"
	"io"
	"net/http"
	"strconv"
	"tec-doc/internal/tec-doc/store/postgres"
	"tec-doc/pkg/errinfo"
	"tec-doc/pkg/model"
)

// @Summary ExcelTemplate
// @Tags excel
// @Description download excel table template
// @ID excel_template
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml
// @Param X-User-Id header string true "ID of user"
// @Param X-Supplier-Id header string true "ID of supplier"
// @Success 200 {array} byte
// @Failure 500 {object} errinfo.errInf
// @Router /api/v1/excel [get]
func (e *externalHttpServer) ExcelTemplate(c *gin.Context) {
	excelTemplate, err := e.service.ExcelTemplateForClient()
	if err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidExcelData))
		return
	}
	c.Data(http.StatusOK, exl.ContentTypeSheetML, excelTemplate)
}

// @Summary ProductsEnrichedExcel
// @Tags excel
// @Description Enrichment excel file, limit entiies in file = 10000
// @ID enrich_excel
// @Produce application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
// @Param excel_file formData file true "excel file"
// Param X-User-Id header string true "ID of user"
// @Param X-Supplier-Id header string true "ID of supplier"
// @Success 200 {array} byte
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/v1/excel/enrichment [post]
func (e *externalHttpServer) GetProductsEnrichedExcel(c *gin.Context) {
	var (
		err      error
		products []model.Product
	)

	if products, err = e.service.LoadFromExcel(c.Request.Body); err != nil {
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
		c.JSON(errinfo.GetErrorInfo(err))
		return
	}
	c.Data(http.StatusOK, exl.ContentTypeSheetML, excel)
}

// @Summary ExcelProductsWithErrors
// @Tags excel
// @Description download excel table template
// @ID excel_products_with_errors
// @Produce json
// @Param InputBody body model.UploadIdRequest true "The input body.<br /> UploadID is ID of previously uploaded task."
// @Param X-User-Id header string true "ID of user"
// @Param X-Supplier-Id header string true "ID of supplier"
// @Success 200 {array} byte
// @Failure 500 {object} errinfo.errInf
// @Router /api/v1/excel/products/errors [post]
func (s *externalHttpServer) ExcelProductsWithErrors(c *gin.Context) {
	var (
		err error
		rq  model.UploadIdRequest
	)
	if err = c.ShouldBindJSON(&rq); err != nil {
		s.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidTaskID))
		return
	}
	var fileRaw []byte
	if fileRaw, err = s.service.ExcelProductsHistoryWithStatus(c, rq.UploadID, postgres.StatusError); err != nil {
		s.logger.Error().Err(err).Msg("can't create excel of products with errors")
		c.JSON(errinfo.GetErrorInfo(errinfo.InternalServerErr))
	}
	c.Data(http.StatusOK, exl.ContentTypeSheetML, fileRaw)
}

// @Summary LoadFromExcel
// @Tags excel
// @Description upload excel table containing products info
// @ID load_from_excel
// @Produce json
// @Param excel_file formData file true "excel file"
// @Param X-User-Id header string true "ID of user"
// @Param X-Supplier-Id header string true "ID of supplier"
// @Success 200 {object} string
// @Failure 400 {object} string
// @Failure 500 {object} string
// @Router /api/v1/excel [post]
func (e *externalHttpServer) LoadFromExcel(c *gin.Context) {
	var (
		supplierID, userID int64
		err                error
	)
	if supplierID, userID, err = CredentialsFromContext(c); err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidTaskID))
		return
	}
	var products []model.Product
	if products, err = e.service.LoadFromExcel(c.Request.Body); err != nil {
		e.logger.Error().Err(err).Send()
		if err.Error() == "empty data" || err == io.EOF {
			c.JSON(errinfo.GetErrorInfo(errinfo.InvalidExcelEmpty))
			return
		}
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidExcelData))
		return
	}

	if err = e.service.AddFromExcel(c, products, supplierID, userID); err != nil {
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
// @Router /api/v1/history/product [post]
func (e *externalHttpServer) GetProductsHistory(c *gin.Context) {
	//_, _, err := CredentialsFromContext(c)
	//if err != nil {
	//	e.logger.Error().Err(err).Send()
	//	c.JSON(errinfo.GetErrorInfo(errinfo.InvalidTaskID))
	//	return
	//}

	var (
		rq  model.UploadIdRequest
		err error
	)
	if err = json.NewDecoder(c.Request.Body).Decode(&rq); err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidTaskID))
		return
	}

	var limit int
	if limit, err = strconv.Atoi(c.Query("limit")); err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidLimit))
		return
	}

	var offset int
	if offset, err = strconv.Atoi(c.Query("offset")); err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidOffset))
		return
	}

	var productsHistory []model.Product
	if productsHistory, err = e.service.GetProductsHistory(c, rq.UploadID, limit, offset); err != nil {
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
// @Success 200 {array} model.TaskPublic
// @Failure 500 {object} errinfo.errInf
// @Router /api/v1/history/task [get]
func (e *externalHttpServer) GetSupplierTaskHistory(c *gin.Context) {
	var (
		supplierID int64
		err        error
	)
	//if supplierID, _, err = CredentialsFromContext(c); err != nil {
	//	e.logger.Error().Err(err).Send()
	//	c.JSON(errinfo.GetErrorInfo(errinfo.InvalidTaskID))
	//	return
	//}

	var limit int
	if limit, err = strconv.Atoi(c.Query("limit")); err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidLimit))
		return
	}

	var offset int
	if offset, err = strconv.Atoi(c.Query("offset")); err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidOffset))
		return
	}
	var rawTasks []model.Task
	if rawTasks, err = e.service.GetSupplierTaskHistory(c, supplierID, limit, offset); err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InternalServerErr))
		return
	}
	var pubTasks = make([]*model.TaskPublic, len(rawTasks))
	for i := range rawTasks {
		pubTasks[i] = &rawTasks[i].TaskPublic
	}
	c.JSON(http.StatusOK, pubTasks)
}

// @Summary GetTecDocArticles
// @Tags product
// @Description to enrichment product by brand and article
// @ID articles_enrichment
// @Accept json
// @Produce json
// @Param limit query string true "limit of contents"
// @Param offset query string true "offset of contents"
// @Param X-User-Id header string true "ID of user"
// @Param X-Supplier-Id header string true "ID of supplier"
// @Param InputBody body model.GetTecDocArticlesRequest true "The input body"
// @Success 200 {array} model.Article
// @Failure 500 {object} errinfo.errInf
// @Router /api/v1/articles/enrichment [post]
func (e *externalHttpServer) GetTecDocArticles(c *gin.Context) {
	//_, _, err := CredentialsFromContext(c)
	//if err != nil {
	//	e.logger.Error().Err(err).Send()
	//	c.JSON(errinfo.GetErrorInfo(errinfo.InvalidTaskID))
	//	return
	//}

	var (
		rq  model.GetTecDocArticlesRequest
		err error
	)
	if err = json.NewDecoder(c.Request.Body).Decode(&rq); err != nil {
		e.logger.Error().Err(err).Send()
		c.JSON(errinfo.GetErrorInfo(errinfo.InvalidTecDocParams))
		return
	}

	var brand *model.Brand
	if brand, err = e.service.GetBrand(rq.Brand); err != nil {
		e.logger.Error().Err(err).Send()
		if errors.Is(err, errinfo.NoTecDocBrandFound) {
			c.JSON(errinfo.GetErrorInfo(err))
		} else {
			c.JSON(errinfo.GetErrorInfo(errinfo.InternalServerErr))
		}
		return
	}

	var articles []model.Article
	if articles, err = e.service.GetArticles(brand.SupplierId, rq.ArticleNumber); err != nil {
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

package externalserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"tec-doc/internal/web/utils"
)

const ContentTypeExcel = "application/vnd.ms-excel"

func (e *externalHttpServer) ExcelTemplate(c *gin.Context) {
	excelTemplate, err := e.service.ExcelTemplateForClient()
	if err != nil {
		utils.JsonError(err, http.StatusInternalServerError, c)
		return
	}
	c.Data(200, ContentTypeExcel, excelTemplate)
}

func (e *externalHttpServer) LoadFromExcel(c *gin.Context) {
	err := e.service.AddFromExcel(c.Request.Body)
	if err != nil {
		utils.JsonError(err, http.StatusInternalServerError, c)
		return
	}
	utils.JsonMessage("success", http.StatusOK, c)
}

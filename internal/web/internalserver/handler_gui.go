package internalserver

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// LoadFromExcelFile	-----> FRONTEND
func (i *internalHttpServer) LoadFromExcelFile(ctx *gin.Context) {
	const TITLE = "Title"
	const ERROR = "Error"
	const MESSAGE = "Message"

	if http.MethodPost == ctx.Request.Method {

		// ---> BACKEND
		i.LoadFromExcel(ctx)
		// TODO узнать как получить результат выполнения бэкэнда

		return
	}

	if http.MethodGet == ctx.Request.Method {
		ctx.HTML(http.StatusOK, "load_excel_file.html", gin.H{
			TITLE:   "Excel Upload",
			MESSAGE: "Загрузить Excel таблицу",
		})
		return
	}
}

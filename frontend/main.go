package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"net/http"
)

var config = &struct {
	BackendAddr  string `json:"backend_addr" default:":8000"`
	FrontendAddr string `json:"frontend_addr" default:":8001"`
}{}

const ContentTypeExcel = "application/vnd.ms-excel"

func init() {
	file, err := ioutil.ReadFile("frontend/config.json")
	if err != nil {
		log.Error().Err(err).Send()
	}
	err = json.Unmarshal(file, config)
	if err != nil {
		log.Error().Err(err).Send()
	}
}

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("./frontend/templates/*.html")

	router.GET("/", indexPage)
	router.POST("/", indexPage)

	server := http.Server{
		Addr:    config.FrontendAddr,
		Handler: router,
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Send()
	}
}

func indexPage(c *gin.Context) {
	if http.MethodGet == c.Request.Method {
		c.HTML(200, "load_excel_file.html", gin.H{
			"redirect": "http://" + config.BackendAddr + "/excel_template",
		})
	}

	if http.MethodPost == c.Request.Method {
		file, _, err := c.Request.FormFile("excel_file")
		if err != nil {
			responseError(err, http.StatusBadRequest, c)
		} else {
			var post *http.Response

			defer file.Close()
			post, err = http.Post("http://"+config.BackendAddr+"/load_from_excel", ContentTypeExcel, file)
			if err != nil {
				responseError(err, http.StatusBadRequest, c)
				return
			}
			if post.StatusCode > 299 {
				data := new(gin.H)
				readAll, _ := ioutil.ReadAll(post.Body)
				err = json.Unmarshal(readAll, data)
				c.HTML(post.StatusCode, "error_excel_file.html", data)
				return
			}
			c.HTML(http.StatusOK, "success.html", gin.H{
				"message":  "Файл успешно загружен",
				"redirect": "/",
			})
		}
	}
}

func responseError(err error, code int, c *gin.Context) {
	c.HTML(code, "error_excel_file.html", gin.H{
		"error":    err.Error(),
		"redirect": "/",
	})
}

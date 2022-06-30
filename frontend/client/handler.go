package client

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"tec-doc/frontend/models"
)

/*
	/excel_template get
	/load_from_excel post
	/task_history get
*/

const (
	frontIndex       = "/"                // POST GET
	frontLoadHistory = "/product_history" // GET

	servExcelTemplate  = "/excel_template"  // GET
	servLoadFromExcel  = "/load_from_excel" // POST
	servProductHistory = "/product_history" // GET

	ContentTypeExcel = "application/vnd.ms-excel"
)

// <<<<<<<<<<<<< Handlers >>>>>>>>>>>>>>

func (cl *Client) downloadsHistory(c *gin.Context) {
	var (
		request  *http.Request
		response *http.Response
		all      []byte
		err      error
		pr       []models.Product
	)
	defer func() {
		if err != nil {
			log.Error().Err(err).Send()
		}
	}()

	request, err = http.NewRequest("GET", cl.createEndpoint(servProductHistory), nil)
	if err != nil {
		responseError(err, http.StatusInternalServerError, c)
		return
	}

	request.Header.Set("limit", strconv.Itoa(cl.limit)) // TODO .......................
	request.Header.Set("offset", "0")                   // TODO .......................

	if response, err = cl.client.Do(request); err != nil || response.StatusCode > 299 {
		responseError(err, response.StatusCode, c)
		return
	}
	if all, err = ioutil.ReadAll(response.Body); err != nil {
		responseError(err, http.StatusInternalServerError, c)
		return
	}
	pr = make([]models.Product, 0, cl.limit)
	if err = json.Unmarshal(all, &pr); err != nil {
		responseError(err, http.StatusInternalServerError, c)
		return
	}
	c.HTML(http.StatusOK, "upload_history.gohtml", pr)
}

func (cl *Client) indexGet(c *gin.Context) {
	c.HTML(200, "load_excel_file.gohtml", gin.H{
		"redirect": cl.createEndpoint(servExcelTemplate),
	})
}

func (cl *Client) indexPost(c *gin.Context) {
	var (
		err  error
		file io.ReadCloser
	)

	file, _, err = c.Request.FormFile("excel_file")
	if err != nil {
		responseError(err, http.StatusBadRequest, c)
		return
	}
	defer func() { _ = file.Close() }()

	response, err := cl.client.Post(cl.createEndpoint(servLoadFromExcel), ContentTypeExcel, file)
	if err != nil {
		responseError(err, http.StatusBadRequest, c)
		return
	}
	if response.StatusCode > 299 {
		var (
			data      *gin.H
			byteSLice []byte
		)
		if byteSLice, err = ioutil.ReadAll(response.Body); err != nil {
			responseError(err, http.StatusInternalServerError, c)
			return
		}
		if err = json.Unmarshal(byteSLice, data); err != nil {
			log.Error().Err(err).Send()
			c.HTML(http.StatusInternalServerError, "error_excel_file.gohtml", data)
			return
		}
		c.HTML(response.StatusCode, "error_excel_file.gohtml", data)
		return
	}
	c.HTML(http.StatusOK, "success.gohtml", gin.H{
		"message":  "Файл успешно загружен",
		"redirect": "/",
	})

}

// <<<<<<<<<<<<< Utils >>>>>>>>>>>>>>

func (cl *Client) createEndpoint(endpoint string) string {
	parse, err := cl.backendURL.Parse(endpoint)
	if err != nil {
		log.Error().Err(err).Send()
		return ""
	}
	return parse.String()
}

func responseError(err error, code int, c *gin.Context) {
	c.HTML(code, "error_excel_file.gohtml", gin.H{
		"error":    err.Error(),
		"redirect": "/",
	})
}

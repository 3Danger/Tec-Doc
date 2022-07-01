package client

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"tec-doc/frontend/models"
)

const (
	frontMainPage      = "/"               // POST GET
	frontExcelTemplate = "/excel_template" // GET

	servExcelTemplate  = "/excel_template"  // GET
	servLoadFromExcel  = "/load_from_excel" // POST
	servProductHistory = "/product_history" // GET
	servTaskHistory    = "/task_history"    //GET

	ContentTypeExcel = "application/vnd.ms-excel"
)

// <<<<<<<<<<<<< Handlers >>>>>>>>>>>>>>

func (cl *Client) indexGet(c *gin.Context) {
	var (
		err   error
		req   *http.Request
		res   *http.Response
		bts   []byte
		tasks []models.Task
	)

	if req, err = cl.createGetRequest(servTaskHistory, c, gin.H{
		"X-User-Id":     "0",
		"X-Supplier-Id": "0",
		"limit":         strconv.Itoa(cl.limit),
		"offset":        "0",
	}); err != nil {
		responseError(err, http.StatusInternalServerError, c)
		return
	}
	if res, err = cl.client.Do(req); err != nil {
		responseError(err, http.StatusInternalServerError, c)
		return
	} else if res.StatusCode > 299 {
		if bts, err = ioutil.ReadAll(res.Body); err != nil {
			responseError(err, http.StatusInternalServerError, c)
			return
		}
		responseError(errors.New(string(bts)), res.StatusCode, c)
		return
	}
	if err = json.NewDecoder(res.Body).Decode(&tasks); err != nil {
		responseError(err, http.StatusInternalServerError, c)
		return
	}

	c.HTML(200, "load_excel_file.gohtml", gin.H{
		"redirect": frontExcelTemplate,
		"tasks":    tasks,
	})
}

func (cl *Client) indexPost(ctx *gin.Context) {
	var (
		err      error
		file     io.ReadCloser
		request  *http.Request
		response *http.Response
	)

	file, _, err = ctx.Request.FormFile("excel_file")
	if err != nil {
		responseError(err, http.StatusBadRequest, ctx)
		return
	}
	defer func() { _ = file.Close() }()

	if request, err = cl.createPostRequest(servLoadFromExcel, file, ctx, gin.H{
		"X-User-Id":     "0",
		"X-Supplier-Id": "0",
	}); err != nil {
		responseError(err, http.StatusInternalServerError, ctx)
		return
	}

	if response, err = cl.client.Do(request); err != nil {
		responseError(err, http.StatusBadRequest, ctx)
		return
	}
	if response.StatusCode > 299 {
		var (
			data      *gin.H
			byteSlice []byte
		)
		if byteSlice, err = ioutil.ReadAll(response.Body); err != nil {
			responseError(err, http.StatusInternalServerError, ctx)
			return
		}
		if err = json.Unmarshal(byteSlice, data); err != nil {
			log.Error().Err(err).Send()
			ctx.HTML(http.StatusInternalServerError, "error_excel_file.gohtml", gin.H{"error": err.Error()})
			return
		}
		ctx.HTML(response.StatusCode, "error_excel_file.gohtml", data)
		return
	}
	ctx.HTML(http.StatusOK, "success.gohtml", gin.H{
		"message":    "Файл успешно загружен",
		"btn_action": "Выгрузить еще",
		"redirect":   frontMainPage,
	})
}

func (cl *Client) downloadExcel(ctx *gin.Context) {
	var (
		err  error
		data []byte
		req  *http.Request
		res  *http.Response
	)
	if req, err = cl.createGetRequest(servExcelTemplate, ctx, gin.H{
		"X-User-Id":     "0",
		"X-Supplier-Id": "0",
	}); err != nil {
		responseError(err, http.StatusInternalServerError, ctx)
		return
	}
	if res, err = cl.client.Do(req); err != nil {
		responseError(err, http.StatusInternalServerError, ctx)
		return
	}
	defer res.Body.Close()
	if data, err = ioutil.ReadAll(res.Body); err != nil {
		responseError(err, http.StatusInternalServerError, ctx)
		return
	}
	if res.StatusCode > 299 {
		responseError(errors.New(string(data)), res.StatusCode, ctx)
		return
	}

	ctx.Data(
		http.StatusOK,
		"Content-Disposition: inline, "+ContentTypeExcel,
		data,
	)
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

func (cl *Client) createPostRequest(endpoint string, body io.Reader, c *gin.Context, h gin.H) (req *http.Request, err error) {
	return createRequest(cl, http.MethodPost, endpoint, body, c, h)
}

func (cl *Client) createGetRequest(endpoint string, c *gin.Context, h gin.H) (req *http.Request, err error) {
	return createRequest(cl, http.MethodGet, endpoint, nil, c, h)
}

func createRequest(cl *Client, method, endpoint string, body io.Reader, c *gin.Context, H gin.H) (req *http.Request, err error) {
	if req, err = http.NewRequest(method, cl.createEndpoint(endpoint), body); err != nil {
		responseError(err, http.StatusInternalServerError, c)
		return
	}
	for k, v := range H {
		req.Header.Set(k, v.(string))
	}
	return req, err
}

package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"tec-doc/frontend/models"
	"tec-doc/internal/model"
)

const (
	ContentTypeExcel = "application/vnd.ms-excel"

	keyUserID     = "X-User-Id"
	keySupplierID = "X-Supplier-Id"
	keyLimit      = "limit"
	keyOffset     = "offset"
)

// <<<<<<<<<<<<< Handlers >>>>>>>>>>>>>>

func (cl *Client) indexGet(c *gin.Context) {
	var (
		err   error
		req   *http.Request
		res   *http.Response
		bts   []byte
		tasks []models.Task
		query = make(url.Values)
	)
	userID, supplierID, limit, offset := getParams(c)
	if req, err = http.NewRequest(http.MethodGet, cl.backendUrlParse("/tasks_history"), nil); err != nil {
		responseError(err, http.StatusInternalServerError, c)
		return
	}
	req.Header.Add(keyUserID, userID)
	req.Header.Add(keySupplierID, supplierID)

	query.Add(keyLimit, limit)
	query.Add(keyOffset, offset)
	req.URL.RawQuery = query.Encode()

	if res, err = cl.client.Do(req); err != nil {
		responseError(err, http.StatusInternalServerError, c)
		return
	}
	if res.StatusCode != 200 {
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

	for i, task := range tasks {
		type ReqStruct struct {
			UploadID int64 `json:"upload_id"`
		}

		rs := ReqStruct{task.ID}
		js, err := json.Marshal(rs)
		if err != nil {
			responseError(err, http.StatusInternalServerError, c)
			return
		}

		reqBodyReader := bytes.NewReader(js)
		req, err = http.NewRequest(http.MethodGet, cl.backendUrlParse("/products_history"), reqBodyReader)

		req.Header.Add(keyUserID, userID)
		req.Header.Add(keySupplierID, supplierID)
		query.Add(keyLimit, limit)
		query.Add(keyOffset, offset)
		req.URL.RawQuery = query.Encode()

		if res, err = cl.client.Do(req); err != nil {
			responseError(err, http.StatusInternalServerError, c)
			return
		}

		if res.StatusCode != 200 {
			if bts, err = ioutil.ReadAll(res.Body); err != nil {
				responseError(err, http.StatusInternalServerError, c)
				return
			}
			responseError(errors.New(string(bts)), res.StatusCode, c)
			return
		}

		products := make([]model.Product, 0)
		if err = json.NewDecoder(res.Body).Decode(&products); err != nil {
			responseError(err, http.StatusInternalServerError, c)
			return
		}
		tasks[i].Products = products
	}

	c.HTML(200, "load_excel_file.gohtml", gin.H{
		"redirect": "/excel_template",
		"tasks":    tasks,
	})
}

func (cl *Client) indexPost(ctx *gin.Context) {
	var (
		err  error
		file io.ReadCloser
		req  *http.Request
		res  *http.Response
	)

	file, _, err = ctx.Request.FormFile("excel_file")
	if err != nil {
		responseError(err, http.StatusBadRequest, ctx)
		return
	}
	defer func() { _ = file.Close() }()

	userID, supplierID, _, _ := getParams(ctx)

	if req, err = http.NewRequest(http.MethodPost, cl.backendUrlParse("/load_from_excel"), file); err != nil {
		responseError(err, http.StatusInternalServerError, ctx)
		return
	}
	req.Header.Set(keyUserID, userID)
	req.Header.Set(keySupplierID, supplierID)

	if res, err = cl.client.Do(req); err != nil {
		responseError(err, http.StatusBadRequest, ctx)
		return
	}
	if res.StatusCode != 200 {
		var (
			data      *gin.H
			byteSlice []byte
		)
		if byteSlice, err = ioutil.ReadAll(res.Body); err != nil {
			responseError(err, http.StatusInternalServerError, ctx)
			return
		}
		if err = json.Unmarshal(byteSlice, data); err != nil {
			log.Error().Err(err).Send()
			ctx.HTML(http.StatusInternalServerError, "error_excel_file.gohtml", gin.H{"error": err.Error()})
			return
		}
		ctx.HTML(res.StatusCode, "error_excel_file.gohtml", data)
		return
	}
	ctx.HTML(http.StatusOK, "success.gohtml", gin.H{
		"message":    "Файл успешно загружен",
		"btn_action": "Выгрузить еще",
		"redirect":   "/",
	})
}

func (cl *Client) downloadExcel(ctx *gin.Context) {
	var (
		err  error
		data []byte
		req  *http.Request
		res  *http.Response
	)
	userID, supplierID, _, _ := getParams(ctx)
	if req, err = http.NewRequest(http.MethodGet, cl.backendUrlParse("/excel_template"), nil); err != nil {
		responseError(err, http.StatusInternalServerError, ctx)
		return
	}
	req.Header.Set(keyUserID, userID)
	req.Header.Set(keySupplierID, supplierID)

	if res, err = cl.client.Do(req); err != nil {
		responseError(err, http.StatusInternalServerError, ctx)
		return
	}
	defer res.Body.Close()
	if data, err = ioutil.ReadAll(res.Body); err != nil {
		responseError(err, http.StatusInternalServerError, ctx)
		return
	}
	if res.StatusCode != 200 {
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

func (cl *Client) backendUrlParse(endpoint string) string {
	parse, err := cl.backendURL.Parse(endpoint)
	if err != nil {
		log.Error().Err(err).Send()
		return ""
	}
	return parse.String()
}

func getParams(c *gin.Context) (userID, supplierID, limit, offset string) {
	w := func(key, defaultValue string) string {
		if c.Request.URL.Query().Has(key) {
			return c.Request.URL.Query().Get(key)
		} else {
			return defaultValue
		}
	}
	userID = w(keyUserID, "0")
	supplierID = w(keySupplierID, "0")
	limit = w(keyLimit, "10")
	offset = w(keyOffset, "0")
	return
}

func responseError(err error, code int, c *gin.Context) {
	c.HTML(code, "error_excel_file.gohtml", gin.H{
		"error":    err.Error(),
		"redirect": "/",
	})
}

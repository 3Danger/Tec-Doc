package tecdoc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"tec-doc/internal/config"
	"tec-doc/internal/model"
	"time"
)

//Client интерфейс с методами для получения запчастей с TecDoc
type TecDocClient interface {
	GetArticles(ctx *context.Context, tecDocCfg config.TecDocConfig, mfrName string) ([]model.Article, error)
	GetBrand(ctx *context.Context, tecDocCfg config.TecDocConfig, mfrName string) (model.Brand, error)
}

type tecDocClient struct {
	http.Client
	baseURL string
}

func NewClient(baseURL string, timeout time.Duration) (*tecDocClient, error) {
	return &tecDocClient{
		Client:  http.Client{Timeout: timeout},
		baseURL: baseURL,
	}, nil
}

func (c *tecDocClient) GetBrand(ctx *context.Context, tecDocCfg config.TecDocConfig, mfrName string) (*model.Brand, error) {
	reqBodyReader := bytes.NewReader([]byte(fmt.Sprintf(
		`{"getBrands":{"articleCountry":"ru", "lang":"ru", "provider":%s}}`, tecDocCfg.ProviderId)))

	req, err := http.NewRequest(http.MethodPost, tecDocCfg.URL, reqBodyReader)
	if err != nil {
		return nil, fmt.Errorf("can't create new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", tecDocCfg.XApiKey)
	//Language code????

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't get response: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read response")
	}

	type respStruct struct {
		Data struct {
			Array []model.Brand `json:"array"`
		} `json:"data"`
		Status int `json:"status"`
	}

	var r respStruct

	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal body: %v", err)
	}

	for _, brand := range r.Data.Array {
		if brand.MfrName == mfrName {
			return &brand, nil
		}
	}

	return nil, fmt.Errorf("no brand found")

}

func (c *tecDocClient) GetArticles(ctx *context.Context, tecDocCfg config.TecDocConfig, dataSupplierId int, article string) ([]model.Article, error) {
	reqBodyReader := bytes.NewReader([]byte(fmt.Sprintf(
		`{
			"getArticles": {
				"articleCountry":"ru", 
				"searchQuery": "%s",
				"searchType": 10,
				"dataSupplierIds": %d,
				"lang":"ru"
			}
		}`, article, dataSupplierId)))

	req, err := http.NewRequest(http.MethodPost, tecDocCfg.URL, reqBodyReader)
	if err != nil {
		return nil, fmt.Errorf("can't create new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", tecDocCfg.XApiKey)
	//Language code????

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't get response: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read response")
	}

	type respStruct struct {
		TotalMatchingArticles int             `json:"totalMatchingArticles"`
		MaxAllowedPage        int             `json:"maxAllowedPage"`
		Articles              []model.Article `json:"articles"`
		Status                int             `json:"status"`
	}

	var r respStruct

	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal body: %v", err)
	}

	if len(r.Articles) == 0 {
		return nil, fmt.Errorf("no articles found")
	}

	return r.Articles, nil
}

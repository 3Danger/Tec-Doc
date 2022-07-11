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

type Client interface {
	GetArticles(ctx context.Context, tecDocCfg config.TecDocConfig, dataSupplierID int, article string) ([]model.Article, error)
	GetBrand(ctx context.Context, tecDocCfg config.TecDocConfig, brandName string) (*model.Brand, error)
}

type tecDocClient struct {
	http.Client
	baseURL string
}

func NewClient(baseURL string, timeout time.Duration) *tecDocClient {
	return &tecDocClient{
		Client:  http.Client{Timeout: timeout},
		baseURL: baseURL,
	}
}

func (c *tecDocClient) GetBrand(ctx context.Context, tecDocCfg config.TecDocConfig, brandName string) (*model.Brand, error) {
	reqBodyReader := bytes.NewReader([]byte(fmt.Sprintf(
		`{"getBrands":{"articleCountry":"ru", "lang":"ru", "provider":%d}}`, tecDocCfg.ProviderId)))

	req, err := http.NewRequest(http.MethodPost, tecDocCfg.URL, reqBodyReader)
	if err != nil {
		return nil, fmt.Errorf("can't create new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", tecDocCfg.XApiKey)

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't get response: %v", err)
	}
	defer resp.Body.Close()

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

	if r.Status != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", r.Status)
	}

	for _, brand := range r.Data.Array {
		if brand.Brand == brandName {
			return &brand, nil
		}
	}

	return nil, fmt.Errorf("no brand found")
}

func (c *tecDocClient) GetArticles(ctx context.Context, tecDocCfg config.TecDocConfig, dataSupplierID int, article string) ([]model.Article, error) {
	reqBodyReader := bytes.NewReader([]byte(fmt.Sprintf(
		`{
			"getArticles": {
				"articleCountry":"ru", 
    			"provider": "%d",
				"searchQuery": "%s",
				"searchType": 10,
				"dataSupplierIds": %d,
				"lang":"ru",
				"includeGenericArticles": true,
				"includeGTINs": true,
				"includeOEMNumbers": true,
				"includeReplacedByArticles": true,
				"includeArticleCriteria": true,
    			"includeImages": true
			}
		}`, tecDocCfg.ProviderId, article, dataSupplierID)))

	req, err := http.NewRequest(http.MethodPost, tecDocCfg.URL, reqBodyReader)
	if err != nil {
		return nil, fmt.Errorf("can't create new request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", tecDocCfg.XApiKey)

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't get response: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("can't read response")
	}

	type respStruct struct {
		TotalMatchingArticles int                `json:"totalMatchingArticles"`
		Articles              []model.ArticleRaw `json:"articles"`
		Status                int                `json:"status"`
	}
	var r respStruct

	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, fmt.Errorf("can't unmarshal body: %v", err)
	}

	if r.Status != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", r.Status)
	}

	if len(r.Articles) == 0 {
		return nil, fmt.Errorf("no articles found")
	}

	return ConvertArticleFromRaw(r.Articles), nil
}

func ConvertArticleFromRaw(rawArticles []model.ArticleRaw) []model.Article {
	articles := make([]model.Article, 0)
	for _, rawArticle := range rawArticles {
		var a model.Article
		a.Brand = rawArticle.MfrName
		a.ArticleNumber = rawArticle.ArticleNumber

		a.ProductGroups = make([]string, 0)
		for _, gr := range rawArticle.GenericArticles {
			a.ProductGroups = append(a.ProductGroups, gr.GenericArticleDescription)
		}

		a.ReplacedByArticles = make([]string, 0)
		for _, rp := range rawArticle.ReplacedByArticles {
			a.ReplacedByArticles = append(a.ReplacedByArticles, rp.ArticleNumber)
		}

		a.Pictures = make([]model.Image, 0)
		a.PanoramicImages = make([]model.Image, 0)
		for _, img := range rawArticle.Images {
			if img.HeaderDescription == "Рисунок" {
				a.Pictures = append(a.Pictures, img)
			} else if img.HeaderDescription == "Панорамное изображение изделия" {
				a.PanoramicImages = append(a.PanoramicImages, img)
			}
		}

		a.EAN = rawArticle.Gtins

		for _, cr := range rawArticle.ArticleCriterias {
			switch cr.CriteriaID {
			case 212:
				a.Weight = cr
			case 1622:
				a.PackageHeight = cr
			case 1621:
				a.PackageWidth = cr
			case 1620:
				a.PackageLength = cr
			}
		}

		a.OEMnumbers = make([]string, 0)
		for _, oem := range rawArticle.OemNumbers {
			a.OEMnumbers = append(a.OEMnumbers, oem.ArticleNumber)
		}

		a.Country = "RU"
		articles = append(articles, a)
	}

	return articles
}

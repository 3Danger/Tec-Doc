package tecdoc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"tec-doc/internal/tec-doc/config"
	"tec-doc/internal/tec-doc/model"
)

type Client interface {
	GetArticles(dataSupplierID int, article string) ([]model.Article, error)
	GetBrand(brandName string) (*model.Brand, error)
}

type tecDocClient struct {
	tecDocCfg config.TecDocClientConfig
	http.Client
	baseURL string
}

func NewClient(baseURL string, tecDocCfg config.TecDocClientConfig) *tecDocClient {
	return &tecDocClient{
		Client:    http.Client{Timeout: tecDocCfg.Timeout},
		baseURL:   baseURL,
		tecDocCfg: tecDocCfg,
	}
}

func (c *tecDocClient) GetBrand(brandName string) (*model.Brand, error) {
	reqBodyReader := bytes.NewReader([]byte(fmt.Sprintf(
		`{"getBrands":{"articleCountry":"ru", "lang":"ru", "provider":%d}}`, c.tecDocCfg.ProviderId)))

	req, err := http.NewRequest(http.MethodPost, c.tecDocCfg.URL, reqBodyReader)
	if err != nil {
		return nil, fmt.Errorf("can't create new request: %v", err)
	}

	req.Header = http.Header{"Content-Type": {"application/json"}, "X-Api-Key": {c.tecDocCfg.XApiKey}}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't get response: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
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
		return nil, fmt.Errorf("can't unmarshal body: %w", err)
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

func (c *tecDocClient) GetArticles(dataSupplierID int, article string) ([]model.Article, error) {
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
		}`, c.tecDocCfg.ProviderId, article, dataSupplierID)))

	req, err := http.NewRequest(http.MethodPost, c.tecDocCfg.URL, reqBodyReader)
	if err != nil {
		return nil, fmt.Errorf("can't create new request: %w", err)
	}

	req.Header = http.Header{"Content-Type": {"application/json"}, "X-Api-Key": {c.tecDocCfg.XApiKey}}

	resp, err := c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("can't get response: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

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
		return nil, fmt.Errorf("can't unmarshal body: %w", err)
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

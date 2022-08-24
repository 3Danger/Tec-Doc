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
				"articleCountry": "RU",
				"provider": 0,
				"searchQuery": "%s",
				"searchType": 0,
				"dataSupplierIds": %d,
				"lang": "ru",
				"perPage": 10,
				"page": 1,
				"includeGenericArticles": true,
				"includeOEMNumbers": true,
				"includeArticleCriteria": true,
				"includeImages": true,
				"assemblyGroupFacetOptions": {"enabled": true, "assemblyGroupType": "P", "includeCompleteTree": false},
				"includeComparableNumbers": true
			}
}`, article, dataSupplierID)))

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

	var r model.TecDocResponse

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

	return ConvertArticleFromRaw(r.Articles, r.AssemblyGroupFacets), nil
}

func ConvertArticleFromRaw(rawArticles []model.ArticleRaw, facets model.AssemblyGroupFacets) []model.Article {
	articles := make([]model.Article, 0)
	for _, rawArticle := range rawArticles {
		var a model.Article

		a.ArticleNumber = rawArticle.ArticleNumber
		a.MfrName = rawArticle.MfrName

		if len(rawArticle.GenericArticles) > 0 {
			a.GenericArticleDescription = rawArticle.GenericArticles[0].GenericArticleDescription

			//TODO
			//legacyID := rawArticle.GenericArticles[0].LegacyArticleID
			//a.LinkageTargets = getLinkageTargets(legacyID)
		}

		for _, oem := range rawArticle.OemNumbers {
			a.OEMnumbers = append(a.OEMnumbers,
				model.OEM{ArticleNumber: oem.ArticleNumber, MfrName: oem.MfrName})
		}

		for _, criteria := range rawArticle.ArticleCriterias {
			convCriteria := convertArticleCriteriaRaw(criteria)
			switch criteria.CriteriaID {
			case 212:
				a.Weight = &convCriteria
			case 1620:
				a.PackageLength = &convCriteria
			case 1621:
				a.PackageWidth = &convCriteria
			case 1622:
				a.PackageHeight = &convCriteria
			case 3653:
				a.PackageDepth = &convCriteria
			default:
				a.ArticleCriteria = append(a.ArticleCriteria, convCriteria)
			}
		}

		for _, img := range rawArticle.Images {
			var imgURL string
			switch {
			case img.ImageURL6400 != "":
				imgURL = img.ImageURL6400
			case img.ImageURL3200 != "":
				imgURL = img.ImageURL3200
			case img.ImageURL1600 == "":
				imgURL = img.ImageURL1600
			case img.ImageURL800 == "":
				imgURL = img.ImageURL800
			case img.ImageURL400 == "":
				imgURL = img.ImageURL400
			case img.ImageURL200 == "":
				imgURL = img.ImageURL200
			case img.ImageURL100 == "":
				imgURL = img.ImageURL100
			case img.ImageURL50 == "":
				imgURL = img.ImageURL50
			}
			a.Images = append(a.Images, imgURL)
		}

		for _, facet := range facets.Counts {
			if facet.Children == 0 {
				a.AssemblyGroupFacets = append(a.AssemblyGroupFacets, facet.AssemblyGroupName)
			}
		}

		articles = append(articles, a)
	}

	return articles
}

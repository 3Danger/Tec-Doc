package tecdoc

import (
	"bytes"
	"fmt"
	"net/http"
	"tec-doc/internal/tec-doc/config"
	"tec-doc/pkg/clients/model"
)

type Client interface {
	GetBrand(brandName string) (*model.Brand, error)
	GetArticles(dataSupplierID int, article string) ([]model.Article, error)
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
	type respStruct struct {
		Data struct {
			Array []model.Brand `json:"array"`
		} `json:"data"`
		Status int `json:"status"`
	}

	var (
		reqBody = []byte(fmt.Sprintf(
			`{"getBrands":{"articleCountry":"ru", "lang":"ru", "provider":%d}}`, c.tecDocCfg.ProviderId))
		resp respStruct
	)

	err := c.doRequest(http.MethodPost, bytes.NewReader(reqBody), &resp)
	if err != nil {
		return nil, err
	}
	if resp.Status != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", resp.Status)
	}

	for _, brand := range resp.Data.Array {
		if brand.Brand == brandName {
			return &brand, nil
		}
	}

	return nil, fmt.Errorf("no brand found")
}

func (c *tecDocClient) GetArticles(dataSupplierID int, article string) ([]model.Article, error) {
	var (
		firstReq = []byte(fmt.Sprintf(
			`{
						  "getArticles": {
							"articleCountry": "RU",
							"searchQuery": "%s",
							"searchType": 0,
							"dataSupplierIds": %d,
							"lang": "ru",
						}
					}`, article, dataSupplierID))

		firstResp = struct {
			TotalMatchingArticles int `json:"totalMatchingArticles"`
			Status                int `json:"status"`
		}{
			0,
			0,
		}
	)

	err := c.doRequest(http.MethodPost, bytes.NewReader(firstReq), &firstResp)
	if err != nil {
		return nil, err
	}

	if firstResp.Status != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", firstResp.Status)
	}

	if firstResp.TotalMatchingArticles == 0 {
		return nil, fmt.Errorf("no articles found")
	}

	const LIMIT = 100
	var (
		stepsNum = firstResp.TotalMatchingArticles/LIMIT + 1
		articles = make([]model.Article, 0)
	)

	for pageNum := 0; pageNum < stepsNum; pageNum++ {
		mainReq := []byte(fmt.Sprintf(
			`{
						"getArticles": {
                                "articleCountry": "RU",
                                "provider": 0,
                                "searchQuery": "%s",
                                "searchType": 0,
                                "dataSupplierIds": %d,
                                "lang": "ru",
                                "perPage": %d,
                                "page": %d,
                                "includeGenericArticles": true,
                                "includeOEMNumbers": true,
                                "includeArticleCriteria": true,
                                "includeImages": true,
                                "assemblyGroupFacetOptions": {"enabled": true, "assemblyGroupType": "P", "includeCompleteTree": false},
                                "includeComparableNumbers": true
                        }
				}`, article, dataSupplierID, LIMIT, pageNum+1))
		var mainResp model.TecDocResponse
		err := c.doRequest(http.MethodPost, bytes.NewReader(mainReq), &mainResp)
		if err != nil {
			return nil, err
		}
		if mainResp.Status != http.StatusOK {
			return nil, fmt.Errorf("request failed with status code: %d", mainResp.Status)
		}
		articles = append(articles, c.ConvertArticleFromRaw(mainResp.Articles, mainResp.AssemblyGroupFacets)...)
	}
	return articles, nil
}

func (c *tecDocClient) ConvertArticleFromRaw(rawArticles []model.ArticleRaw, facets model.AssemblyGroupFacets) []model.Article {
	articles := make([]model.Article, 0)
	for _, rawArticle := range rawArticles {
		var a model.Article

		a.ArticleNumber = rawArticle.ArticleNumber
		a.MfrName = rawArticle.MfrName

		if len(rawArticle.GenericArticles) > 0 {
			a.GenericArticleDescription = rawArticle.GenericArticles[0].GenericArticleDescription
			legacyID := rawArticle.GenericArticles[0].LegacyArticleID
			a.LinkageTargets, _ = c.Applicability(legacyID)
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

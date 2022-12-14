package tecdoc

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog"
	"net/http"
	"sync"
	"tec-doc/internal/tec-doc/config"
	"tec-doc/internal/tec-doc/store/postgres"
	"tec-doc/pkg/errinfo"
	"tec-doc/pkg/model"
	"time"
)

type Client interface {
	GetBrand(brandName string) (*model.Brand, error)
	GetArticles(dataSupplierID int, article string) ([]model.Article, error)
	Enrichment(products []model.Product) (productsEnriched []model.ProductEnriched)
	ConvertToCharacteristics(pe *model.ProductEnriched) *model.ProductCharacteristics
}

type tecDocClient struct {
	tecDocCfg config.TecDocClientConfig
	http.Client
	baseURL string
	logger  *zerolog.Logger
}

func NewClient(baseURL string, tecDocCfg config.TecDocClientConfig, log *zerolog.Logger) Client {
	return &tecDocClient{
		Client:    http.Client{Timeout: tecDocCfg.Timeout},
		baseURL:   baseURL,
		tecDocCfg: tecDocCfg,
		logger:    log,
	}
}

type getBrandType struct {
	Data struct {
		Array []model.Brand `json:"array"`
	} `json:"data"`
	Status int `json:"status"`

	//Что бы не делать миллион одинаковых запросов на каждый товар
	time       time.Time
	providerId int
}

var respBrand getBrandType

func (c *tecDocClient) GetBrand(brandName string) (*model.Brand, error) {
	var (
		reqBody = []byte(fmt.Sprintf(
			`{"getBrands":{"articleCountry":"ru", "lang":"ru", "provider":%d}}`, c.tecDocCfg.ProviderId))
		//resp respStruct
	)
	if time.Now().Sub(respBrand.time) > (time.Minute*5) ||
		respBrand.providerId != c.tecDocCfg.ProviderId ||
		respBrand.Status != http.StatusOK {

		respBrand.providerId = c.tecDocCfg.ProviderId
		respBrand.time = time.Now()
		err := c.doRequest(http.MethodPost, bytes.NewReader(reqBody), &respBrand)
		if err != nil {
			return nil, fmt.Errorf("failed to do request: %w", err)
		}
		if respBrand.Status != http.StatusOK {
			return nil, fmt.Errorf("request failed with status code: %d", respBrand.Status)
		}
	}
	for _, brand := range respBrand.Data.Array {
		if brand.Brand == brandName {
			return &brand, nil
		}
	}

	return nil, errinfo.NoTecDocBrandFound
}

func (c *tecDocClient) GetArticles(dataSupplierID int, article string) ([]model.Article, error) {
	var (
		firstReq = []byte(fmt.Sprintf(
			`{
						  "getArticles": {
							"articleCountry": "RU",
							"provider": %d,
							"searchQuery": "%s",
							"searchType": 0,
							"dataSupplierIds": %d,
							"lang": "ru",
						}
					}`, c.tecDocCfg.ProviderId, article, dataSupplierID))

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
		return nil, errinfo.NoTecDocArticlesFound
	}

	if firstResp.TotalMatchingArticles > 1 {
		return nil, errinfo.MoreThanOneArticlesFound
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
                                "provider": %d,
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
								"assemblyGroupFacetOptions": {
                                    "enabled": true,
                                    "assemblyGroupType": "P",
                                    "includeCompleteTree": true
                                }
                        }
				}`, c.tecDocCfg.ProviderId, article, dataSupplierID, LIMIT, pageNum+1))
		var mainResp model.TecDocResponse
		if err = c.doRequest(http.MethodPost, bytes.NewReader(mainReq), &mainResp); err != nil {
			return nil, err
		}
		if mainResp.Status != http.StatusOK {
			return nil, fmt.Errorf("request failed with status code: %d", mainResp.Status)
		}
		articles = append(articles, c.ConvertArticleFromRaw(mainResp)...)
	}
	return articles, nil
}

func (c *tecDocClient) ConvertArticleFromRaw(mainResp model.TecDocResponse) []model.Article {
	rawArticles := mainResp.Articles
	articles := make([]model.Article, 0)
	for _, rawArticle := range rawArticles {
		var (
			a   model.Article
			err error
		)

		a.ArticleNumber = rawArticle.ArticleNumber
		a.MfrName = rawArticle.MfrName
		a.GenericArticleDescription = rawArticle.GenericArticles[0].GenericArticleDescription

		if a.CrossNumbers, err = c.GetCrossNumbers(a.ArticleNumber); err != nil {
			c.logger.Error().Err(err).Send()
		}

		if len(rawArticle.GenericArticles) > 0 {
			if a.LinkageTargets, err = c.Applicability(rawArticle.GenericArticles[0].LegacyArticleID); err != nil {
				c.logger.Error().Err(err).Send()
			}
		}

		for _, oem := range rawArticle.OemNumbers {
			a.OEMnumbers = append(a.OEMnumbers,
				model.OEM{ArticleNumber: oem.ArticleNumber, MfrName: oem.MfrName})
		}

		for _, criteria := range rawArticle.ArticleCriterias {
			convCriteria := convertArticleCriteriaRaw(criteria)
			if convCriteria == nil {
				continue
			}
			if criteria.CriteriaID == 212 || criteria.CriteriaID == 1620 || criteria.CriteriaID == 1621 ||
				criteria.CriteriaID == 1622 || criteria.CriteriaID == 1623 || criteria.CriteriaID == 3653 {
				a.PackageArticleCriteria = append(a.PackageArticleCriteria, *convCriteria)
			} else {
				a.ArticleCriteria = append(a.ArticleCriteria, *convCriteria)
			}
		}

		for _, img := range rawArticle.Images {
			var imgURL string
			switch {
			case img.ImageURL3200 != "":
				imgURL = img.ImageURL3200
			case img.ImageURL1600 != "":
				imgURL = img.ImageURL1600
			case img.ImageURL800 != "":
				imgURL = img.ImageURL800
			case img.ImageURL400 != "":
				imgURL = img.ImageURL400
			case img.ImageURL200 != "":
				imgURL = img.ImageURL200
			case img.ImageURL100 != "":
				imgURL = img.ImageURL100
			case img.ImageURL50 != "":
				imgURL = img.ImageURL50
			}
			a.Images = append(a.Images, imgURL)
		}

		a.AssemblyGroupName = getAssemblyGroupName(mainResp.AssemblyGroupFacets)
		articles = append(articles, a)
	}
	return articles
}

func (c *tecDocClient) GetCrossNumbers(articleNumber string) ([]model.CrossNumbers, error) {

	type firstRespStruct struct {
		TotalMatchingArticles int `json:"totalMatchingArticles"`
		Status                int `json:"status"`
	}

	var (
		firstReq = []byte(fmt.Sprintf(
			`{
						"getArticles": {
							"articleCountry": "RU",
        					"provider": %d,
							"searchQuery": "%s",
							"searchType": 3,
							"lang": "ru",
						}
					}`, c.tecDocCfg.ProviderId, articleNumber))
		firstResp firstRespStruct
	)

	err := c.doRequest(http.MethodPost, bytes.NewReader(firstReq), &firstResp)
	if err != nil {
		return nil, fmt.Errorf("failed to do first GetCrossNumbers request: %w", err)
	}

	if firstResp.Status != http.StatusOK {
		return nil, fmt.Errorf("etCrossNumbers request failed with status code: %d", firstResp.Status)
	}

	if firstResp.TotalMatchingArticles == 0 {
		return nil, fmt.Errorf("no comparable numbers found")
	}

	type mainRespStruct struct {
		Articles []struct {
			ArticleNumber string `json:"articleNumber"`
			MfrName       string `json:"mfrName"`
		} `json:"articles"`
		Status int `json:"status"`
	}

	const LIMIT = 100
	var (
		stepsNum        = firstResp.TotalMatchingArticles/LIMIT + 2
		mainResp        mainRespStruct
		crossNumbers    = make([]model.CrossNumbers, 0)
		replaceArticles = make([]struct {
			ArticleNumber string `json:"articleNumber"`
			MfrName       string `json:"mfrName"`
		}, 0)
	)

	for pageNum := 1; pageNum < stepsNum; pageNum++ {
		mainReq := []byte(fmt.Sprintf(
			`{
						"getArticles": {
							"articleCountry": "RU",
        					"provider": %d,
							"searchQuery": "%s",
							"searchType": 3,
							"perPage": %d,
							"page":	%d,
							"lang": "ru",
						}
					}`, c.tecDocCfg.ProviderId, articleNumber, LIMIT, pageNum))
		err := c.doRequest(http.MethodPost, bytes.NewReader(mainReq), &mainResp)
		if err != nil {
			return nil, fmt.Errorf("failed to do main GetCrossNumbers request: %w", err)
		}

		if mainResp.Status != http.StatusOK {
			return nil, fmt.Errorf("main GetCrossNumbers request failed with status code: %d", mainResp.Status)
		}

		replaceArticles = append(replaceArticles, mainResp.Articles...)
	}

	for _, replaceArticle := range replaceArticles {
		if replaceArticle.ArticleNumber != articleNumber {
			crossNumbers = append(crossNumbers, model.CrossNumbers{
				ArticleNumber: replaceArticle.ArticleNumber,
				MfrName:       replaceArticle.MfrName})
		}
	}
	return crossNumbers, nil
}

func (c *tecDocClient) Enrichment(products []model.Product) []model.ProductEnriched {
	const GOROUTINESMAX = 1000
	var (
		productsEnriched = make([]model.ProductEnriched, len(products), len(products))
		wg               sync.WaitGroup
		stepsNum         = len(products)/GOROUTINESMAX + 1
	)

	//TODO переработать
	for step := 0; step < stepsNum; step++ {
		for i := step * GOROUTINESMAX; i < (step+1)*GOROUTINESMAX && i < len(products); i++ {
			wg.Add(1)
			go func(product *model.Product, toEnrich *model.ProductEnriched, wg *sync.WaitGroup) {
				defer wg.Done()
				err := c.SingleEnrichment(product, toEnrich)
				if err != nil {
					c.logger.Error().Str("tecDocClient", "Enrichment").Err(err).Send()
					toEnrich.Status = postgres.StatusError
					_, toEnrich.ErrorResponse = errinfo.GetErrorInfo(err)
				}
			}(&products[i], &productsEnriched[i], &wg)
		}
		wg.Wait()
	}
	return productsEnriched
}

func (c *tecDocClient) SingleEnrichment(product *model.Product, enriched *model.ProductEnriched) error {
	var (
		brand    *model.Brand
		articles []model.Article
	)
	*enriched = model.ProductEnriched{Product: *product}
	brand, err := c.GetBrand(product.Brand)
	if err != nil {
		return err
	}
	if articles, err = c.GetArticles(brand.SupplierId, product.Article); err != nil {
		return err
	}
	(*enriched).Article = articles[0]
	return nil
}

func getAssemblyGroupName(facets model.AssemblyGroupFacets) string {
	for j := facets.Total - 1; j >= 0; j-- {
		if facets.Counts[j].Children == 0 {
			return facets.Counts[j].AssemblyGroupName
		}
	}

	return ""
}

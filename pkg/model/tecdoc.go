package model

import (
	"time"
)

type UploadIdRequest struct {
	UploadID int64 `json:"uploadID" form:"uploadID" example:"1"`
}

type GetTecDocArticlesRequest struct {
	Brand         string `json:"brand"`
	ArticleNumber string `json:"articleNumber"`
}

type Brand struct {
	SupplierId int    `json:"dataSupplierId"`
	Brand      string `json:"mfrName"`
}

type TecDocResponse struct {
	TotalMatchingArticles int          `json:"totalMatchingArticles"`
	Articles              []ArticleRaw `json:"articles"`
	Status                int          `json:"status"`
}

type ArticleRaw struct {
	MfrName         string `json:"mfrName"`
	ArticleNumber   string `json:"articleNumber"`
	GenericArticles []struct {
		GenericArticleID          int    `json:"genericArticleId"`
		GenericArticleDescription string `json:"genericArticleDescription"`
		LegacyArticleID           int    `json:"legacyArticleId"`
	} `json:"genericArticles"`
	OemNumbers []struct {
		ArticleNumber      string `json:"articleNumber"`
		MfrID              int    `json:"mfrId"`
		MfrName            string `json:"mfrName"`
		MatchesSearchQuery bool   `json:"matchesSearchQuery"`
	} `json:"oemNumbers"`
	ArticleCriterias   []ArticleCriteriaRaw `json:"articleCriteria"`
	Images             []Image              `json:"images"`
	SearchQueryMatches []struct {
		MatchType   string `json:"matchType"`
		Description string `json:"description"`
		Match       string `json:"match"`
	} `json:"searchQueryMatches"`
}

type Article struct {
	ArticleNumber             string            `json:"articleNumber,omitempty"`
	MfrName                   string            `json:"mfrName,omitempty"`
	GenericArticleDescription string            `json:"genericArticleDescription,omitempty"`
	OEMnumbers                []OEM             `json:"oemNumbers,omitempty"`
	CrossNumbers              []CrossNumbers    `json:"crossNumbers,omitempty"`
	ArticleCriteria           []ArticleCriteria `json:"articleCriteria,omitempty"`
	PackageArticleCriteria    []ArticleCriteria `json:"packageArticleCriteria"`
	LinkageTargets            []LinkageTargets  `json:"linkageTargets,omitempty"`
	Images                    []string          `json:"images,omitempty"`
}

type ArticleCriteriaRaw struct {
	CriteriaID              int    `json:"criteriaId,omitempty"`
	CriteriaDescription     string `json:"criteriaDescription"`
	CriteriaAbbrDescription string `json:"criteriaAbbrDescription"`
	CriteriaUnitDescription string `json:"criteriaUnitDescription,omitempty"`
	CriteriaType            string `json:"criteriaType"`
	RawValue                string `json:"rawValue"`
	FormattedValue          string `json:"formattedValue"`
	ImmediateDisplay        bool   `json:"immediateDisplay"`
	IsMandatory             bool   `json:"isMandatory"`
	IsInterval              bool   `json:"isInterval"`
}

type ArticleCriteria struct {
	CriteriaDescription     string `json:"criteriaDescription"`
	CriteriaAbbrDescription string `json:"criteriaAbbrDescription"`
	CriteriaUnitDescription string `json:"criteriaUnitDescription,omitempty"`
	CriteriaType            string `json:"criteriaType"`
	RawValue                string `json:"rawValue"`
	FormattedValue          string `json:"formattedValue"`
}

type Image struct {
	ImageURL50        string `json:"imageURL50"`
	ImageURL100       string `json:"imageURL100"`
	ImageURL200       string `json:"imageURL200"`
	ImageURL400       string `json:"imageURL400"`
	ImageURL800       string `json:"imageURL800"`
	ImageURL1600      string `json:"imageURL1600"`
	ImageURL3200      string `json:"imageURL3200"`
	ImageURL6400      string `json:"imageURL6400"`
	FileName          string `json:"fileName"`
	TypeDescription   string `json:"typeDescription"`
	TypeKey           int    `json:"typeKey"`
	HeaderDescription string `json:"headerDescription"`
	HeaderKey         int    `json:"headerKey"`
}

type OEM struct {
	ArticleNumber string `json:"articleNumber"`
	MfrName       string `json:"mfrName"`
}

type CrossNumbers struct {
	ArticleNumber string `json:"articleNumber"`
	MfrName       string `json:"mfrName"`
}

type TaskPublic struct {
	ID                int64     `json:"id"`
	UploadDate        time.Time `json:"uploadDate"`
	Status            int       `json:"status"`
	ProductsProcessed int       `json:"productsProcessed"`
	ProductsFailed    int       `json:"productsFailed"`
	ProductsTotal     int       `json:"productsTotal"`
}

type Task struct {
	SupplierID int64     `json:"supplierID"`
	UpdateDate time.Time `json:"updateDate"`
	IP         string    `json:"ip"`
	UserID     int64     `json:"userID"`
	TaskPublic
}

type ProductEnriched struct {
	Product
	Article
}

type Product struct {
	ID              int64     `json:"id"`
	UploadID        int64     `json:"uploadId"`
	Article         string    `json:"article"`
	ArticleSupplier string    `json:"articleSupplier"`
	Brand           string    `json:"brand"`
	Barcode         string    `json:"barcode"`
	Subject         string    `json:"subject"`
	Price           int       `json:"price"`
	UploadDate      time.Time `json:"uploadDate"`
	UpdateDate      time.Time `json:"updateDate"`
	Status          int       `json:"status"`
	ErrorResponse   string    `json:"errorResponse"`
	Amount          int       `json:"amount" default:"1"`
}

type ProductCharacteristics struct {
	Object          string  `json:"Предмет"`
	Brand           string  `json:"Бренд"`
	Subject         string  `json:"Категория"`
	ArticleSupplier string  `json:"Артикул товара"`
	Article         string  `json:"Артикул производителя"`
	Barcode         string  `json:"Штрихкод товара"`
	Price           int     `json:"Розничная цена, в руб"`
	GenArticleDescr string  `json:"Наименование"`
	OEMnumbers      string  `json:"ОЕМ номер"`
	Weight          float64 `json:"Вес с упаковкой (кг)"`
	Height          float64 `json:"Высота упаковки"`
	Depth           float64 `json:"Глубина упаковки"`
	Width           float64 `json:"Ширина упаковки"`
	Description     string  `json:"Описание"`
	Targets         string  `json:"Марка автомобиля"`
	Photo           string  `json:"Фото"`
	Amount          int     `json:"Комплектация"`
	ErrorResponse   string  `json:"Ошибки"`
}

type LinkageTargets struct {
	LinkageTargetId        int    `json:"linkageTargetId"`
	MfrName                string `json:"mfrName"`
	VehicleModelSeriesName string `json:"vehicleModelSeriesName"`
	BeginYearMonth         string `json:"beginYearMonth"`
	EndYearMonth           string `json:"endYearMonth"`
}

type (
	// GetLinkageTargetsResponse для запроса
	GetLinkageTargets struct {
		PerPage              int              `json:"perPage"`
		Page                 int              `json:"page"`
		LinkageTargetCountry string           `json:"linkageTargetCountry"`
		Lang                 string           `json:"lang"`
		LinkageTargetIds     []map[string]any `json:"linkageTargetIds"`
	}
	GetLinkageTargetsResponse struct {
		GetLinkageTargets GetLinkageTargets `json:"getLinkageTargets"`
	}
)

// Data для записи ответа первого запроса
type Data struct {
	Data struct {
		Array []ArticleLinkages `json:"array"`
	} `json:"data"`
	Status int `json:"status"`
}

type ArticleLinkages struct {
	ArticleLinkages struct {
		LinkingTargetId []struct {
			LinkingTargetId int `json:"linkingTargetId"`
		} `json:"array"`
	} `json:"articleLinkages"`
}

// LinkageTargetsResponse для записи результата
type LinkageTargetsResponse struct {
	Total          int              `json:"total"`
	LinkageTargets []LinkageTargets `json:"linkageTargets"`
	Status         int              `json:"status"`
}

package model

import (
	"time"
)

type Brand struct {
	SupplierId int    `json:"dataSupplierId"`
	Brand      string `json:"mfrName"`
}

type TecDocResponse struct {
	TotalMatchingArticles int                 `json:"totalMatchingArticles"`
	Articles              []ArticleRaw        `json:"articles"`
	Status                int                 `json:"status"`
	AssemblyGroupFacets   AssemblyGroupFacets `json:"assemblyGroupFacets"`
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
	ComparableNumbers  []interface{}        `json:"comparableNumbers"`
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
	ArticleCriteria           []ArticleCriteria `json:"articleCriteria,omitempty"`
	Weight                    *ArticleCriteria  `json:"weight,omitempty"`
	PackageLength             *ArticleCriteria  `json:"packageLength,omitempty"`
	PackageWidth              *ArticleCriteria  `json:"packageWidth,omitempty"`
	PackageHeight             *ArticleCriteria  `json:"packageHeight,omitempty"`
	PackageDepth              *ArticleCriteria  `json:"packageDepth,omitempty"`
	LinkageTargets            []LinkageTargets  `json:"linkageTargets,omitempty"`
	Images                    []string          `json:"images,omitempty"`
	ComparableNumbers         []interface{}     `json:"comparableNumbers,omitempty"`
	AssemblyGroupFacets       []string          `json:"assemblyGroupFacets,omitempty"`
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

type AssemblyGroupFacets struct {
	Total  int `json:"total"`
	Counts []struct {
		AssemblyGroupNodeID int    `json:"assemblyGroupNodeId"`
		AssemblyGroupName   string `json:"assemblyGroupName"`
		AssemblyGroupType   string `json:"assemblyGroupType"`
		Children            int    `json:"children,omitempty"`
		Count               int    `json:"count"`
		ParentNodeID        int    `json:"parentNodeId,omitempty"`
	} `json:"counts"`
}

type Task struct {
	ID                int64     `json:"id"`
	SupplierID        int64     `json:"supplierID"`
	UserID            int64     `json:"userID"`
	UploadDate        time.Time `json:"uploadDate"`
	UpdateDate        time.Time `json:"updateDate"`
	IP                string    `json:"ip"`
	Status            int       `json:"status"`
	ProductsProcessed int       `json:"productsProcessed"`
	ProductsFailed    int       `json:"productsFailed"`
	ProductsTotal     int       `json:"productsTotal"`
}

type Product struct {
	ID                  int64     `json:"id"`
	UploadID            int64     `json:"uploadId"`
	CardNumber          int       `json:"cardNumber"`
	Article             string    `json:"article"`
	ProviderArticle     string    `json:"providerArticle"`
	ManufacturerArticle string    `json:"manufacturerArticle"`
	Brand               string    `json:"brand"`
	SKU                 string    `json:"sku"`
	Category            string    `json:"category"`
	Price               int       `json:"price"`
	UploadDate          time.Time `json:"uploadDate"`
	UpdateDate          time.Time `json:"updateDate"`
	Status              int       `json:"status"`
	ErrorResponse       string    `json:"errorResponse"`
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

package model

import "time"

type Brand struct {
	SupplierId int    `json:"dataSupplierId"`
	Brand      string `json:"mfrName"`
}

type ArticleRaw struct {
	MfrName       string `json:"mfrName"`
	ArticleNumber string `json:"articleNumber"`

	GenericArticles []struct {
		GenericArticleDescription string `json:"genericArticleDescription"`
	} `json:"genericArticles"`

	Gtins []string `json:"gtins"`

	OemNumbers []struct {
		ArticleNumber      string `json:"articleNumber"`
		MfrID              int    `json:"mfrId"`
		MfrName            string `json:"mfrName"`
		MatchesSearchQuery bool   `json:"matchesSearchQuery"`
	} `json:"oemNumbers"`

	ReplacedByArticles []struct {
		ArticleNumber  string `json:"articleNumber"`
		DataSupplierID int    `json:"dataSupplierId"`
		MfrID          int    `json:"mfrId"`
		MfrName        string `json:"mfrName"`
	} `json:"replacedByArticles"`

	ArticleCriterias []ArticleCriteria `json:"articleCriteria"`

	Images []Image `json:"images"`

	DataSupplierID int `json:"dataSupplierId"`
	MfrID          int `json:"mfrId"`
}

type Article struct {
	Brand              string
	ArticleNumber      string
	ProductGroups      []string
	ReplacedByArticles []string
	Pictures           []Image
	PanoramicImages    []Image
	EAN                []string
	Weight             ArticleCriteria
	PackageHeight      ArticleCriteria
	PackageWidth       ArticleCriteria
	PackageLength      ArticleCriteria
	OEMnumbers         []string
	RelatedVehicles    []string
	Country            string
}

type ArticleCriteria struct {
	CriteriaID              int    `json:"criteriaId"`
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

type Image struct {
	ImageURL50        string `json:"imageURL50"`
	ImageURL100       string `json:"imageURL100"`
	ImageURL200       string `json:"imageURL200"`
	ImageURL400       string `json:"imageURL400"`
	ImageURL800       string `json:"imageURL800"`
	ImageURL1600      string `json:"imageURL1600"`
	ImageURL3200      string `json:"imageURL3200"`
	FileName          string `json:"fileName"`
	TypeDescription   string `json:"typeDescription"`
	TypeKey           int    `json:"typeKey"`
	HeaderDescription string `json:"headerDescription"`
	HeaderKey         int    `json:"headerKey"`
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
	Article             string    `json:"article"`
	CardNumber          int       `json:"cardNumber"`
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

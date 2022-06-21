package model

type Brand struct {
	DataSupplierId int    `json:"dataSupplierId"`
	MfrName        string `json:"mfrName"`
}

type Article struct {
	DataSupplierID     int          `json:"dataSupplierId"`
	ArticleNumber      string       `json:"articleNumber"`
	MfrID              int          `json:"mfrId"`
	MfrName            string       `json:"mfrName"`
	SearchQueryMatches []QueryMatch `json:"searchQueryMatches"`
}

type QueryMatch struct {
	MatchType   string `json:"matchType"`
	Description string `json:"description"`
	Match       string `json:"match"`
	MfrID       int    `json:"mfrId,omitempty"`
	MfrName     string `json:"mfrName,omitempty"`
}

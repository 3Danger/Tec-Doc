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

type Task struct {
	ID          int64
	SupplierID  int64
	UserID      int64
	Description string
}

type Product struct {
	ID            int64
	UploadID      int64
	Article       string
	Brand         string
	Status        int
	ErrorResponse string
	Description   string
}

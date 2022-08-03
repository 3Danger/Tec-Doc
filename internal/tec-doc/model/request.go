package model

type GetProductsHistoryRequest struct {
	UploadID int64 `json:"UploadID" example:"1"`
}

type GetTecDocArticlesRequest struct {
	Brand         string `json:"Brand" example:"BOSCH"`
	ArticleNumber string `json:"ArticleNumber" example:"0451103274"`
}

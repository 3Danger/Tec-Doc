package model

type GetProductsHistoryRequest struct {
	UploadID int64 `json:"UploadID" example:"1"`
}

type GetTecDocArticlesRequest struct {
	Brand         string `json:"Brand"`
	ArticleNumber string `json:"ArticleNumber"`
}

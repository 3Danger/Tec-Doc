package models

import "time"

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

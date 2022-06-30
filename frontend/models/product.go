package models

import "time"

type Product struct {
	ID                  int64     `json:"id"`
	UploadID            int64     `json:"upload_id"`
	Article             string    `json:"article"`
	CardNumber          int       `json:"card_number"`
	ProviderArticle     string    `json:"provider_article"`
	ManufacturerArticle string    `json:"manufacturer_article"`
	Brand               string    `json:"brand"`
	SKU                 string    `json:"sku"`
	Category            string    `json:"category"`
	Price               int       `json:"price"`
	UploadDate          time.Time `json:"upload_date"`
	UpdateDate          time.Time `json:"update_date"`
	Status              int       `json:"status"`
	ErrorResponse       string    `json:"error_response"`
}

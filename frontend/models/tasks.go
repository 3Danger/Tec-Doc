package models

import "time"

type Task struct {
	ID                int64     `json:"id"`
	SupplierID        int64     `json:"supplier_id"`
	UserID            int64     `json:"user_id"`
	UploadDate        time.Time `json:"upload_date"`
	UpdateDate        time.Time `json:"update_date"`
	IP                string    `json:"ip"`
	Status            int       `json:"status"`
	ProductsProcessed int       `json:"products_processed"`
	ProductsFailed    int       `json:"products_failed"`
	ProductsTotal     int       `json:"products_total"`
}

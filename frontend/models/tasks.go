package models

import (
	"tec-doc/internal/model"
	"time"
)

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
	Products          []model.Product
}

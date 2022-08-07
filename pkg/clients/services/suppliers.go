package services

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

//go:generate tg client -go --services . --outPath ../suppliers

// @tg jsonRPC-server log metrics
type Suppliers interface {
	GetOldSupplierID(ctx context.Context, supplierID uuid.UUID) (oldSupplierID int, err error)
}

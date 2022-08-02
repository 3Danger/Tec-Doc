package service

import (
	"context"
	"tec-doc/internal/tec-doc/model"
)

func (s *Service) GetSupplierTaskHistory(ctx context.Context, supplierID int64, limit int, offset int) ([]model.Task, error) {
	return s.database.GetSupplierTaskHistory(ctx, nil, supplierID, limit, offset)
}

func (s *Service) GetProductsHistory(ctx context.Context, uploadID int64, limit int, offset int) ([]model.Product, error) {
	return s.database.GetProductsHistory(ctx, nil, uploadID, limit, offset)
}

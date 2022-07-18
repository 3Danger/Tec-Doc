package service

import (
	"context"
	"tec-doc/internal/tec-doc/model"
	"tec-doc/internal/tec-doc/store/postgres"
)

func (s *Service) GetSupplierTaskHistory(ctx context.Context, tx postgres.Transaction, supplierID int64, limit int, offset int) ([]model.Task, error) {
	return s.database.GetSupplierTaskHistory(ctx, tx, supplierID, limit, offset)
}

func (s *Service) GetProductsHistory(ctx context.Context, tx postgres.Transaction, uploadID int64, limit int, offset int) ([]model.Product, error) {
	return s.database.GetProductsHistory(ctx, tx, uploadID, limit, offset)
}

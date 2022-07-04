package service

import (
	"context"
	"tec-doc/internal/model"
	"tec-doc/internal/store/postgres"
)

func (s *Service) GetSupplierTaskHistory(ctx context.Context, tx postgres.Transaction, supplierID int64, limit int, offset int) ([]model.Task, error) {
	return s.database.GetSupplierTaskHistory(ctx, tx, supplierID, limit, offset)
}

func (s *Service) GetProductsHistory(ctx context.Context, tx postgres.Transaction, uploadID int64, limit int, offset int) ([]model.Product, error) {
	return s.database.GetProductsHistory(ctx, tx, uploadID, limit, offset)
}

////GetProductHistory in ctx must be upload_id, limit and offset
//func (s *Service) GetProductHistory(ctx context.Context) ([]model.Product, error) {
//	uploadId := ctx.Value("upload_id").(int64)
//	limit := ctx.Value("limit").(int)
//	offset := ctx.Value("offset").(int)
//	return s.database.GetProductsHistory(ctx, nil, uploadId, limit, offset)
//}

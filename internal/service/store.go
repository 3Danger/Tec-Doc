package service

import (
	"context"
	"tec-doc/internal/model"
)

func (s *Service) GetSupplierTaskHistory(ctx context.Context, supplierID int64, limit int, offset int) ([]model.Task, error) {
	return s.database.GetSupplierTaskHistory(ctx, nil, supplierID, limit, offset)
}

//GetProductHistory in ctx must be upload_id, limit and offset
func (s *Service) GetProductHistory(ctx context.Context) ([]model.Product, error) {
	uploadId := ctx.Value("upload_id").(int64)
	limit := ctx.Value("limit").(int)
	offset := ctx.Value("offset").(int)
	return s.database.GetProductsHistory(ctx, nil, uploadId, limit, offset)
}

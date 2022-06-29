package service

import (
	"context"
	"tec-doc/internal/model"
)

func (s *Service) GetSupplierTaskHistory(ctx context.Context, supplierID int64, limit int, offset int) ([]model.Task, error) {
	return s.database.GetSupplierTaskHistory(ctx, supplierID, limit, offset)
}

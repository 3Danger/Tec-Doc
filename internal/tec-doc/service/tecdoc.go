package service

import (
	"context"
	"tec-doc/internal/tec-doc/model"
)

func (s *Service) GetArticles(ctx context.Context, dataSupplierID int, article string) ([]model.Article, error) {
	return s.tecDocClient.GetArticles(ctx, s.conf.TecDoc, dataSupplierID, article)
}

func (s *Service) GetBrand(ctx context.Context, brandName string) (*model.Brand, error) {
	return s.tecDocClient.GetBrand(ctx, s.conf.TecDoc, brandName)
}
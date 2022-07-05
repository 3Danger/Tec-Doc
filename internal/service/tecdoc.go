package service

import (
	"context"
	"tec-doc/internal/config"
	"tec-doc/internal/model"
)

func (s *Service) GetArticles(ctx context.Context, tecDocCfg config.TecDocConfig, dataSupplierID int, article string) ([]model.Article, error) {
	return s.tecDocClient.GetArticles(ctx, tecDocCfg, dataSupplierID, article)
}

func (s *Service) GetBrand(ctx context.Context, tecDocCfg config.TecDocConfig, brandName string) (*model.Brand, error) {
	return s.tecDocClient.GetBrand(ctx, tecDocCfg, brandName)
}

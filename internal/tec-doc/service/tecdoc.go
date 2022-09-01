package service

import (
	"tec-doc/pkg/model"
)

func (s *Service) GetArticles(dataSupplierID int, article string) ([]model.Article, error) {
	return s.tecDocClient.GetArticles(dataSupplierID, article)
}

func (s *Service) GetBrand(brandName string) (*model.Brand, error) {
	return s.tecDocClient.GetBrand(brandName)
}

func (s *Service) Enrichment(product []model.Product) ([]model.ProductEnriched, error) {
	return s.tecDocClient.Enrichment(product)
}

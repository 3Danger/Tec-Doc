package service

import (
	"tec-doc/internal/tec-doc/model"
)

func (s *Service) GetArticles(dataSupplierID int, article string) ([]model.Article, error) {
	return s.tecDocClient.GetArticles(dataSupplierID, article)
}

func (s *Service) GetBrand(brandName string) (*model.Brand, error) {
	return s.tecDocClient.GetBrand(brandName)
}

package service

import (
	"tec-doc/pkg/model"
)

func (s *Service) GetArticles(dataSupplierID int, article string) ([]model.Article, error) {
	articles, err := s.tecDocClient.GetArticles(dataSupplierID, article)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (s *Service) GetBrand(brandName string) (*model.Brand, error) {
	brand, err := s.tecDocClient.GetBrand(brandName)
	if err != nil {
		return nil, err
	}
	return brand, nil
}

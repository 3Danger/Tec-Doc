package service

import (
	"tec-doc/internal/tec-doc/model"
	"tec-doc/pkg/errinfo"
)

func (s *Service) GetArticles(dataSupplierID int, article string) ([]model.Article, error) {
	articles, err := s.tecDocClient.GetArticles(dataSupplierID, article)
	if err != nil && err.Error() == "no articles found" {
		return nil, errinfo.NoTecDocArticlesFound
	}
	return articles, err
}

func (s *Service) GetBrand(brandName string) (*model.Brand, error) {
	brand, err := s.tecDocClient.GetBrand(brandName)
	if err != nil && err.Error() == "no brand found" {
		return nil, errinfo.NoTecDocBrandFound
	}
	return brand, err
}

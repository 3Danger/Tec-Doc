package tecdoc

import (
	"math"
	"sort"
	"strconv"
	"strings"
	"tec-doc/pkg/model"
)

func (t *tecDocClient) ConvertToCharacteristics(pe *model.ProductEnriched) *model.ProductCharacteristics {
	var width, height, depth, weight float64
	if pe.PackageArticleCriteria != nil {
		width, height, depth, weight = t.Sizes(pe.Article.PackageArticleCriteria)
	}
	if width == 0 {
		width, height, depth, weight = t.Sizes(pe.Article.ArticleCriteria)
	}
	result := &model.ProductCharacteristics{
		Object:          "Автозапчасти",
		Brand:           pe.Product.Brand,
		Subject:         pe.Product.Subject,
		ArticleSupplier: pe.Product.ArticleSupplier,
		Article:         pe.Product.Article,
		Barcode:         pe.Product.Barcode,
		Price:           pe.Product.Price,
		GenArticleDescr: pe.GenericArticleDescription,
		OEMnumbers:      t.OemCross(pe.OEMnumbers, pe.CrossNumbers),
		Weight:          weight,
		Height:          height,
		Depth:           depth,
		Width:           width,
		Description:     t.ArticleCriteria(pe.Article.ArticleCriteria),
		Targets:         t.LinkageTargets(pe.Article.LinkageTargets),
		Photo:           strings.Join(pe.Images, ";"),
		Amount:          pe.Product.Amount,
		ErrorResponse:   pe.ErrorResponse,
	}
	return result
}

func (*tecDocClient) LinkageTargets(lts []model.LinkageTargets) string {
	resultMap := make(map[string]struct{})
	for i := range lts {
		resultMap[lts[i].MfrName] = struct{}{}
	}
	result := make([]string, 0, len(resultMap))
	for k := range resultMap {
		result = append(result, k)
	}
	sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
	return strings.Join(result, ";")
}

func (*tecDocClient) ArticleCriteria(ac []model.ArticleCriteria) string {
	result := make([]string, len(ac))
	for i := range result {
		result[i] = ac[i].CriteriaAbbrDescription + " " + ac[i].RawValue + ac[i].CriteriaUnitDescription
	}
	return strings.Join(result, ";")
}

func (*tecDocClient) OemCross(oems []model.OEM, cross []model.CrossNumbers) string {
	mapBuff := make(map[string]interface{})
	for i := range oems {
		mapBuff[oems[i].ArticleNumber] = struct{}{}
	}
	for i := range cross {
		mapBuff[cross[i].ArticleNumber] = struct{}{}
	}

	result := make([]string, 0, len(mapBuff))
	for k := range mapBuff {
		result = append(result, k)
	}
	return strings.Join(result, ";")
}
func normalise(rawValue, criteriaUnit string) float64 {
	var res float64
	rawValue = strings.Replace(rawValue, ",", ".", 1)
	res, _ = strconv.ParseFloat(rawValue, 64)
	if res == 0 {
		return 0
	}
	switch strings.ToLower(criteriaUnit) {
	case "мм", "mm":
		res /= 10
	case "гр":
		res /= 1000
	}
	return res
}
func (*tecDocClient) Sizes(pe []model.ArticleCriteria) (width, height, depth, weight float64) {
	for i := range pe {
		criteria := strings.ToLower(pe[i].CriteriaAbbrDescription)
		if strings.Contains(criteria, "вес") {
			weight = normalise(pe[i].RawValue, pe[i].CriteriaUnitDescription)
		} else if strings.Contains(criteria, "длина") {
			depth = normalise(pe[i].RawValue, pe[i].CriteriaUnitDescription)
		} else if strings.Contains(criteria, "ширина") {
			width = normalise(pe[i].RawValue, pe[i].CriteriaUnitDescription)
		} else if strings.Contains(criteria, "высота") {
			height = normalise(pe[i].RawValue, pe[i].CriteriaUnitDescription)
		}
	}
	max := math.Max(math.Max(width, height), depth)
	if width == 0 {
		width = max
	}
	if height == 0 {
		height = max
	}
	if depth == 0 {
		depth = max
	}
	return width, height, depth, weight
}

//func handleCriteriaUnitDescriotnion(ac *model.ArticleCriteria)

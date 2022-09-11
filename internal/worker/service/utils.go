package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
	"tec-doc/pkg/model"
	"unicode"
)

func (s *service) makeUploadBody(pe *model.ProductEnriched) (body io.Reader, err error) {
	type M map[string]interface{}
	var bodyMap = make(M)
	width, height, depth, weight := s.Sizes(&pe.Article)
	characteristics := []interface{}{
		M{"Предмет": "Автозапчасти"},
		M{"Бренд": pe.Product.Brand},
		M{"Категория": pe.Product.Subject},
		M{"Артикул товара": pe.Product.ArticleSupplier},
		M{"Артикул производителя": pe.Product.Article},
		M{"Штрихкод товара": pe.Product.Barcode},
		M{"Розничная цена, в руб": pe.Product.Price},
		M{"Наименование": pe.GenericArticleDescription},
		M{"ОЕМ номер": s.OemCross(pe.Article.OEMnumbers, pe.Article.CrossNumbers)},
		M{"Вес с упаковкой (кг)": weight},
		M{"Высота упаковки": height},
		M{"Глубина упаковки": depth},
		M{"Ширина упаковки": width},
		M{"Описание": s.ArticleCriteria(pe.Article.ArticleCriteria)},
		M{"Марка автомобиля": s.LinkageTargets(pe.Article.LinkageTargets)},
		M{"Фото": strings.Join(pe.Images, ";")},
		M{"Комплектация": fmt.Sprintf("%dшт.", pe.Set)},
	}
	bodyMap["vendorCode"] = pe.Product.ArticleSupplier
	bodyMap["characteristics"] = characteristics
	buff := new(bytes.Buffer)
	if err = json.NewEncoder(buff).Encode(bodyMap); err != nil {
		return nil, err
	}
	return buff, nil
}

func (*service) LinkageTargets(lts []model.LinkageTargets) string {
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

func (*service) ArticleCriteria(ac []model.ArticleCriteria) string {
	result := make([]string, len(ac))
	for i := range result {
		result[i] = ac[i].CriteriaAbbrDescription + " " + ac[i].RawValue + ac[i].CriteriaUnitDescription
	}
	return strings.Join(result, ";")
}

func (*service) OemCross(oems []model.OEM, cross []model.CrossNumbers) string {
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
	if criteriaUnit == "мм" {
		res *= 1000
	}
	return res
}
func (s *service) Sizes(pe *model.Article) (width, height, depth, weight float64) {
	pack := pe.PackageArticleCriteria
	if pack != nil {
		for i := range pack {
			switch pack[i].CriteriaAbbrDescription {
			case "Вес":
				weight = normalise(pack[i].RawValue, "")
			case "Длина упаковки":
				depth = normalise(pack[i].RawValue, "")
			case "Ширина упаковки":
				width = normalise(pack[i].RawValue, "")
			case "Высота упаковки":
				height = normalise(pack[i].RawValue, "")
			}
		}
	}
	max := math.Max(math.Max(width, height), depth)
	if max == 0 {
		art := pe.ArticleCriteria
		for i := range art {
			abbr := strings.ToLower(art[i].CriteriaAbbrDescription)
			if strings.Contains(abbr, "размеры") {
				x := strings.TrimFunc(art[i].RawValue, func(r rune) bool { return unicode.IsNumber(r) })
				sizes := strings.Split(art[i].RawValue, x)
				switch len(sizes) {
				case 3:
					width = normalise(sizes[2], art[i].CriteriaUnitDescription)
					height = normalise(sizes[1], art[i].CriteriaUnitDescription)
					depth = normalise(sizes[0], art[i].CriteriaUnitDescription)
				case 2:
					width = normalise(sizes[1], art[i].CriteriaUnitDescription)
					height = normalise(sizes[0], art[i].CriteriaUnitDescription)
				case 1:
					width = normalise(sizes[0], art[i].CriteriaUnitDescription)
				}
				break
			}
		}
	}
	if width == 0 {
		width = max
	}
	if height == 0 {
		height = max
	}
	if depth == 0 {
		height = 0
	}
	return width, height, depth, weight
}

//func handleCriteriaUnitDescriotnion(ac *model.ArticleCriteria)

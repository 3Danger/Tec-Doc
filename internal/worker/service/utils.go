package service

import (
	"bytes"
	"encoding/json"
	"io"
	"tec-doc/pkg/model"
)

func (s *service) makeUploadBody(productEnriched *model.ProductEnriched) (body io.Reader, err error) {
	ch := s.enricher.ConvertToCharacteristics(productEnriched)
	type M map[string]interface{}
	var bodyMap = make(M)
	characteristics := []interface{}{
		M{"Предмет": ch.Object},
		M{"Бренд": ch.Brand},
		M{"Категория": ch.Subject},
		M{"Артикул товара": ch.ArticleSupplier},
		M{"Артикул производителя": ch.Article},
		M{"Штрихкод товара": ch.Barcode},
		M{"Розничная цена, в руб": ch.Price},
		M{"Наименование": ch.GenArticleDescr},
		M{"ОЕМ номер": ch.OEMnumbers},
		M{"Вес с упаковкой (кг)": ch.Weight},
		M{"Высота упаковки": ch.Height},
		M{"Глубина упаковки": ch.Depth},
		M{"Ширина упаковки": ch.Width},
		M{"Описание": ch.Description},
		M{"Марка автомобиля": ch.Targets},
		M{"Фото": ch.Photo},
		M{"Комплектация": ch.Amount},
	}
	bodyMap["vendorCode"] = ch.ArticleSupplier
	bodyMap["characteristics"] = characteristics
	buff := new(bytes.Buffer)
	if err = json.NewEncoder(buff).Encode(bodyMap); err != nil {
		return nil, err
	}
	return buff, nil
}

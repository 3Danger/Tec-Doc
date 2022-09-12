package service

import (
	"bytes"
	"encoding/json"
	"io"
	"tec-doc/pkg/model"
)

func (s *service) makeUploadBody(pe *model.ProductCharacteristics) (body io.Reader, err error) {
	type M map[string]interface{}
	var bodyMap = make(M)
	characteristics := []interface{}{
		M{"Предмет": pe.Object},
		M{"Бренд": pe.Brand},
		M{"Категория": pe.Subject},
		M{"Артикул товара": pe.ArticleSupplier},
		M{"Артикул производителя": pe.Article},
		M{"Штрихкод товара": pe.Barcode},
		M{"Розничная цена, в руб": pe.Price},
		M{"Наименование": pe.GenArticleDescr},
		M{"ОЕМ номер": pe.OEMnumbers},
		M{"Вес с упаковкой (кг)": pe.Weight},
		M{"Высота упаковки": pe.Height},
		M{"Глубина упаковки": pe.Depth},
		M{"Ширина упаковки": pe.Width},
		M{"Описание": pe.Description},
		M{"Марка автомобиля": pe.Targets},
		M{"Фото": pe.Photo},
		M{"Комплектация": pe.Set},
	}
	bodyMap["vendorCode"] = pe.ArticleSupplier
	bodyMap["characteristics"] = characteristics
	buff := new(bytes.Buffer)
	if err = json.NewEncoder(buff).Encode(bodyMap); err != nil {
		return nil, err
	}
	return buff, nil
}

package model

type CreateProductCardRequest struct {
	JsonRPCVersion      string          `json:"jsonrpc"`
	ID                  int64           `json:"id"`
	ImtID               int64           `json:"imtID"`      //номер карточки?
	VendorCode          string          `json:"vendorCode"` //артикул поставщика?
	ManufacturerArticle string          `json:"Артикул производителя"`
	Characteristics     Characteristics `json:"characteristics"`
}

type Characteristics struct {
	Brand           string   `json:"Бренд"`
	SKU             []string `json:"SKU"`
	ProductCategory string   `json:"Категория"`
	Price           int      `json:"Цена"` // нет в списке характеристик котента

	//из tec-doc
	ArticleNumber      string   `json:"Номер артикула"`  // нет в списке характеристик котента
	ProductGroups      []string `json:"Товарные группы"` // нет в списке характеристик котента
	ReplacedByArticles []string `json:"Замена"`          // нет в списке характеристик котента
	Image              struct{} `json:"Рисунок"`         // нет документации по добавлению
	//панорамное изображение изделия?
	EAN             string   `json:"EAN"`             // (штрихкод )нет в списке характеристик котента
	Weight          string   `json:"Вес (кг)"`        // единицы измерения
	PackageHeight   string   `json:"Высота упаковки"` // единицы измерения
	PackageWidth    string   `json:"Ширина упаковки"` // единицы измерения
	PackageLength   string   `json:"Длина упаковки"`  // единицы измерения
	OEM             []string `json:"ОЕМ номер"`       //в базе кириллица в аббревиатуре
	RelatedVehicles []string // нет в списке характеристик котента и текдока
	Country         string   //страна чего? нет в списке характеристик котента
}

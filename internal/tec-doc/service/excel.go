package service

import (
	"github.com/gin-gonic/gin"
	exl "github.com/xuri/excelize/v2"
	"tec-doc/pkg/model"
	"time"
)

var styleExcelHeader = &exl.Style{
	Fill: exl.Fill{},
	Font: &exl.Font{
		Bold:   true,
		Family: "Fira Sans Book",
		Size:   8,
		Color:  "c21f6b"},
	Lang: "ru",
}

var styleExcel = &exl.Style{
	Fill: exl.Fill{},
	Font: &exl.Font{
		Family: "Fira Sans Book",
		Size:   7,
		Color:  "731a6f"},
	Lang: "ru",
}

func (s *Service) ExcelTemplateForClient() ([]byte, error) {
	f := exl.NewFile()
	defer func() { _ = f.Close() }()
	nameSheet := "Products"
	f.SetSheetName(f.GetSheetName(0), nameSheet)

	// Set styles & length
	styleHeaderId, err := f.NewStyle(styleExcelHeader)
	if err != nil {
		return nil, err
	}
	styleId, err := f.NewStyle(styleExcel)
	if err != nil {
		return nil, err
	}
	_ = f.SetRowStyle(nameSheet, 1, 2, styleHeaderId)
	_ = f.SetRowStyle(nameSheet, 2, 1000, styleId)
	_ = f.SetRowHeight(nameSheet, 1, 20)
	_ = f.SetColWidth(nameSheet, "A", "A", 9)
	_ = f.SetColWidth(nameSheet, "B", "C", 42)
	_ = f.SetColWidth(nameSheet, "D", "D", 15)
	_ = f.SetColWidth(nameSheet, "E", "E", 20)

	// Set values
	_ = f.SetCellValue(nameSheet, "A1", "Бренд")
	_ = f.SetCellValue(nameSheet, "B1", "Артикул поставщика (уникальный артикул)")
	_ = f.SetCellValue(nameSheet, "C1", "Артикул производителя (артикул tec-doc)")
	_ = f.SetCellValue(nameSheet, "D1", "Цена товара")
	_ = f.SetCellValue(nameSheet, "E1", "Штрих-код")
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (s *Service) AddFromExcel(ctx *gin.Context, products []model.Product, supplierID int64, userID int64) error {
	tx, err := s.database.Transaction(ctx)
	if err != nil {
		return err
	}

	uploaderId, err := s.database.CreateTask(ctx, tx, supplierID, userID, ctx.ClientIP(), time.Now().UTC())
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	for i := 0; i < len(products); i++ {
		products[i].UploadID = uploaderId
	}
	if err = s.database.SaveIntoBuffer(ctx, tx, products); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	if err = tx.Commit(ctx); err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return nil
}

func (s *Service) GetProductsEnrichedExcel(productsPoor []model.Product) (data []byte, err error) {
	products, err := s.tecDocClient.Enrichment(productsPoor)
	if err != nil {
		return nil, err
	}
	f := exl.NewFile()
	defer func() { _ = f.Close() }()
	nameSheet := "Details about products"
	f.SetSheetName(f.GetSheetName(0), nameSheet)

	// Set styles & length
	styleHeaderId, err := f.NewStyle(styleExcelHeader)
	if err != nil {
		return nil, err
	}
	styleId, err := f.NewStyle(styleExcel)
	if err != nil {
		return nil, err
	}
	{
		// Set header into excel file
		//TODO узнать в чем дело, почему первый стиль перекрывает собой другой стиль
		_ = f.SetRowStyle(nameSheet, 1, 2, styleHeaderId)
		_ = f.SetRowStyle(nameSheet, 2, 10, styleId)
		setWidths := func(width []float64) (err error) {
			for i := 0; i < len(width) && err == nil; i++ {
				a, b := rune('A'+i), rune('A'+i+1)
				err = f.SetColWidth(nameSheet, string(a), string(b), width[i])
			}
			return
		}
		if err = setWidths([]float64{14, 15, 25, 25, 18, 18, 40, 40, 40, 40, 40, 40}); err != nil {
			return nil, err
		}
		setValue := func(x, y int, value interface{}) {
			axis, _ := exl.CoordinatesToCellName(x, y)
			_ = f.SetCellValue(nameSheet, axis, value)
		}

		// Add values from struct Products
		setValue(1, 1, "Бренд")
		setValue(2, 1, "Категория")
		setValue(3, 1, "Артикул поставщика (уникальный артикул)")
		setValue(4, 1, "Артикул производителя (артикул tec-doc)")
		setValue(5, 1, "Штрих-код")
		setValue(6, 1, "Цена товара")

		// Add values from struct Articles
		setValue(7, 1, "Описание")
		setValue(8, 1, "Cross numbers")
		setValue(9, 1, "Размерность")
		setValue(10, 1, "ArticleCriterias")
		setValue(11, 1, "PackageArticleCriterias")
		setValue(12, 1, "Применимости")
		setValue(13, 1, "Изображения")

		for i, p := range products {
			i += 2
			setValue(1, i, p.Product.Brand)
			setValue(2, i, p.Product.Subject)
			setValue(3, i, p.Product.ArticleSupplier)
			setValue(4, i, p.Product.Article)
			setValue(5, i, p.Product.Barcode)
			setValue(6, i, p.Product.Price)

			setValue(7, i, p.Article.GenericArticleDescription)
			setValue(8, i, p.Article.OEMnumbers)
			setValue(9, i, p.Article.CrossNumbers)
			setValue(10, i, p.Article.ArticleCriteria)
			setValue(11, i, p.Article.PackageArticleCriteria)
			setValue(12, i, p.Article.LinkageTargets)
			setValue(13, i, p.Article.Images)
		}
		err = f.Save()

	}
	//A B C D E F G H I K L M N O P Q R S T V X Y Z
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

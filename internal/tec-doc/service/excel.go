package service

import (
	"bytes"
	"fmt"
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
	var productsEnriched []model.ProductEnriched

	if productsEnriched, err = s.tecDocClient.Enrichment(productsPoor); err != nil {
		return nil, err
	}

	f := exl.NewFile()
	defer func() { _ = f.Close() }()
	f.SetSheetName(f.GetSheetName(0), "Details about products")
	sw, err := f.NewStreamWriter(f.GetSheetName(0))

	// Set width for columns
	for i, w := range []float64{14, 15, 25, 25, 18, 13, 25, 50, 40, 40, 8, 40, 40} {
		i++
		if err = sw.SetColWidth(i, i, w); err != nil {
			return nil, err
		}
	}
	if err = sw.SetRow("A1", []interface{}{
		"Бренд",
		"Категория",
		"Артикул поставщика (уникальный артикул)",
		"Артикул производителя (артикул tec-doc)",
		"Штрих-код",
		"Цена товара",
		"Описание",
		"Cross numbers",
		"Размерность",
		"ArticleCriterias",
		"PackageArticleCriterias",
		"Применимости",
		"Изображения",
	}, exl.RowOpts{Height: 15},
	); err != nil {
		return nil, err
	}

	for i, p := range productsEnriched {
		axis := fmt.Sprintf("A%d", i+2)
		err = sw.SetRow(axis, []interface{}{
			p.Product.Brand,
			p.Product.Subject,
			p.Product.ArticleSupplier,
			p.Product.Article,
			p.Product.Barcode,
			p.Product.Price,
			p.Article.GenericArticleDescription,
			p.Article.OEMnumbers,
			p.Article.CrossNumbers,
			p.Article.ArticleCriteria,
			p.Article.PackageArticleCriteria,
			p.Article.LinkageTargets,
			p.Article.Images},
			exl.RowOpts{Height: 15})
		if err != nil {
			return nil, err
		}
	}
	if err = sw.Flush(); err != nil {
		return nil, err
	}
	var buffer *bytes.Buffer
	if buffer, err = f.WriteToBuffer(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

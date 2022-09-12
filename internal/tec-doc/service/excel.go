package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	exl "github.com/xuri/excelize/v2"
	"io"
	"strconv"
	"strings"
	"tec-doc/pkg/model"
	"time"
	"unicode"
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
	Alignment: &exl.Alignment{
		Vertical: "center",
		WrapText: true,
	},
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
	_ = f.SetCellValue(nameSheet, "F1", "Комплектация")
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (e *Service) LoadFromExcel(bodyData io.Reader) (products []model.Product, err error) {
	var rows [][]string
	f, err := exl.OpenReader(bodyData)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	list := f.GetSheetList()
	if len(list) == 0 {
		return nil, errors.New("empty data")
	}
	rows, err = f.GetRows("Products")
	if len(rows) < 2 {
		return nil, errors.New("empty data")
	}
	products = make([]model.Product, len(rows[1:]))
	for i := range products {
		if err = e.parseExcelRow(&products[i], rows[i+1]); err != nil {
			return nil, err
		}
	}
	return products, nil
}

func (e *Service) parseExcelRow(p *model.Product, row []string) (err error) {
	if len(row) < 5 {
		return errors.New("row is invalid")
	}
	p.Brand = row[0]
	p.ArticleSupplier = row[1]
	p.Article = row[2]
	if p.Price, err = strconv.Atoi(row[3]); err != nil {
		return err
	}
	p.Barcode = row[4]
	if len(row) >= 6 && len(row[5]) != 0 {
		n := strings.LastIndexFunc(row[5], func(r rune) bool { return unicode.IsNumber(r) })
		if p.Amount, err = strconv.Atoi(row[5][:n+1]); err != nil {
			return err
		}
	}
	return nil
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
	style, err := f.NewStyle(styleExcel)
	if err != nil {
		return nil, err
	}
	f.SetSheetName(f.GetSheetName(0), "Details about products")
	sw, err := f.NewStreamWriter(f.GetSheetName(0))

	// Set width for columns
	for i, w := range []float64{20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20, 20} {
		i++
		if err = sw.SetColWidth(i, i, w); err != nil {
			return nil, err
		}
	}
	if err = sw.SetRow("A1", []interface{}{
		"Предмет",
		"Бренд",
		"Категория",
		"Артикул товара",
		"Артикул производителя",
		"Штрихкод товара",
		"Розничная цена, в руб",
		"Наименование",
		"ОЕМ номер",
		"Вес с упаковкой (кг)",
		"Высота упаковки",
		"Глубина упаковки",
		"Ширина упаковки",
		"Описание",
		"Марка автомобиля",
		"Фото",
		"Комплектация",
		"Ошибки",
	}, exl.RowOpts{Height: 15},
	); err != nil {
		return nil, err
	}

	for i, p := range productsEnriched {
		ch := s.tecDocClient.ConvertToCharacteristics(&p)
		axis := fmt.Sprintf("A%d", i+2)
		err = sw.SetRow(axis, []interface{}{
			ch.Object,
			ch.Brand,
			ch.Subject,
			ch.ArticleSupplier,
			ch.Article,
			ch.Barcode,
			ch.Price,
			ch.GenArticleDescr,
			ch.OEMnumbers,
			ch.Weight,
			ch.Height,
			ch.Depth,
			ch.Width,
			ch.Description,
			ch.Targets,
			ch.Photo,
			ch.Amount,
			ch.ErrorResponse,
		},
			exl.RowOpts{Height: 15, StyleID: style})
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

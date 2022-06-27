package service

import (
	"errors"
	"github.com/gin-gonic/gin"
	exl "github.com/xuri/excelize/v2"
	"io"
	"strconv"
	"tec-doc/internal/model"
	"time"
)

var styleExcelHeader = &exl.Style{
	Fill: exl.Fill{},
	Font: &exl.Font{
		Bold:   true,
		Family: "Fira Sans Book",
		Size:   14,
		Color:  "c21f6b"},
	Lang: "ru",
}

var styleExcel = &exl.Style{
	Fill: exl.Fill{},
	Font: &exl.Font{
		Family: "Fira Sans Book",
		Size:   13,
		Color:  "731a6f"},
	Lang: "ru",
}

func (s *Service) ExcelTemplateForClient() ([]byte, error) {
	f := exl.NewFile()
	defer func() { _ = f.Close() }()
	nameSheet := "Products"
	index := f.NewSheet(nameSheet)
	// Set values
	_ = f.SetCellValue(nameSheet, "A1", "Номер карточки")
	_ = f.SetCellValue(nameSheet, "B1", "Артикул поставщика (уникальный артикул)")
	_ = f.SetCellValue(nameSheet, "C1", "Артикул производителя")
	_ = f.SetCellValue(nameSheet, "D1", "Бренд")
	_ = f.SetCellValue(nameSheet, "E1", "SKU")
	_ = f.SetCellValue(nameSheet, "F1", "Категория товара")
	_ = f.SetCellValue(nameSheet, "G1", "Цена товара")

	// Set styles & length
	styleHeaderId, err := f.NewStyle(styleExcelHeader)
	if err != nil {
		return nil, err
	}
	styleId, err := f.NewStyle(styleExcel)
	if err != nil {
		return nil, err
	}
	_ = f.SetRowStyle(nameSheet, 1, 1, styleHeaderId)
	_ = f.SetRowStyle(nameSheet, 2, 1000, styleId)
	_ = f.SetRowHeight(nameSheet, 1, 20)
	_ = f.SetColWidth(nameSheet, "A", "F", 24)
	_ = f.SetColWidth(nameSheet, "G", "G", 19)
	_ = f.SetColWidth(nameSheet, "D", "E", 9)
	_ = f.SetColWidth(nameSheet, "B", "B", 45)
	f.SetActiveSheet(index)
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (s *Service) loadFromExcel(bodyData io.Reader) (products []model.Product, err error) {
	f, err := exl.OpenReader(bodyData)
	if err != nil {
		return nil, err
	}
	list := f.GetSheetList()
	if len(list) == 0 {
		return nil, errors.New("empty data")
	}
	rows, err := f.GetRows(list[0])
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, errors.New("empty data")
	}
	products = make([]model.Product, len(rows[1:]))
	for i := range products {
		if err = parseExcelRow(&products[i], rows[i+1]); err != nil {
			return nil, err
		}
	}
	return products, nil
}

func parseExcelRow(p *model.Product, row []string) (err error) {
	if len(row) < 7 {
		return errors.New("row is invalid")
	}
	if p.NumberOfCard, err = strconv.Atoi(row[0]); err != nil {
		return err
	}
	p.ProviderArticle = row[1]
	p.ManufactureArticle = row[2]
	p.Brand = row[3]
	p.SKU = row[4]
	p.Category = row[5]
	if p.Price, err = strconv.Atoi(row[6]); err != nil {
		return err
	}
	return nil
}

func (s *Service) AddFromExcel(bodyData io.Reader, ctx *gin.Context) error {
	products, err := s.loadFromExcel(bodyData)
	if err != nil {
		return err
	}

	// TODO TRANSACTION!!
	//tx, err := s.database.Transaction()
	//if err != nil {
	//	return err
	//}
	//defer func() { _ = tx.Rollback() }()

	//TODO don't forget take fields from ctx *gin.Context,
	// a supplierID and userID,
	// that getting from middleware:
	// Authorize(next *gin.Context)
	uploaderId, err := s.database.CreateTask(ctx, 1, 1, "sd", time.Now().UTC())
	if err != nil {
		return err
	}
	for i := 0; i < len(products); i++ {
		products[i].UploadID = uploaderId
	}
	if err = s.database.SaveIntoBuffer(ctx, products); err != nil {
		return err
	}
	// TODO TX COMMIT()!!
	return nil
}

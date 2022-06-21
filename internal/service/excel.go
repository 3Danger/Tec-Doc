package service

import (
	"bytes"
	"errors"
	exl "github.com/xuri/excelize/v2"
	"strconv"
)

type Product struct {
	NumberOfCard       int
	ProviderArticle    string
	ManufactureArticle string
	Brand              string
	SKU                string
	ProductCategory    string
	ProductPrice       float64
}

func parseExcelRow(p *Product, row []string) (err error) {
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
	p.ProductCategory = row[5]
	if p.ProductPrice, err = strconv.ParseFloat(row[6], 64); err != nil {
		return err
	}
	return nil
}

func (s *Service) ExcelTemplateForClient(nameSheet string) ([]byte, error) {
	f := exl.NewFile()
	defer func() { _ = f.Close() }()
	if nameSheet == "" {
		nameSheet = "Layer1"
	}
	index := f.NewSheet(nameSheet)
	_ = f.SetCellValue(nameSheet, "A1", "Номер карточки")
	_ = f.SetCellValue(nameSheet, "B1", "Артикул поставщика (уникальный артикул)")
	_ = f.SetCellValue(nameSheet, "C1", "Артикул производителя")
	_ = f.SetCellValue(nameSheet, "D1", "Бренд")
	_ = f.SetCellValue(nameSheet, "E1", "SKU")
	_ = f.SetCellValue(nameSheet, "F1", "Категория товара")
	_ = f.SetCellValue(nameSheet, "G1", "Цена товара")
	f.SetActiveSheet(index)
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (s *Service) LoadFromExcel(rawData []byte) (products []Product, err error) {
	f, err := exl.OpenReader(bytes.NewReader(rawData))
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
	products = make([]Product, len(rows))
	for i := range products {
		err = parseExcelRow(&products[i], rows[i])
		if err != nil {
			return nil, err
		}
	}
	return products, nil
}

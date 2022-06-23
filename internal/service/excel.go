package service

import (
	"database/sql"
	"errors"
	exl "github.com/xuri/excelize/v2"
	"io"
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

func (s *Service) loadFromExcel(bodyData io.Reader) (products []*Product, err error) {
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
	products = make([]*Product, len(rows[1:]))
	for i := range products {
		products[i] = new(Product)
		if err = parseExcelRow(products[i], rows[i+1]); err != nil {
			return nil, err
		}
	}
	return products, nil
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

func (s *Service) AddFromExcel(bodyData io.Reader) error {
	products, err := s.loadFromExcel(bodyData)
	if err != nil {
		return err
	}
	//TODO will be implement it when be initialized database
	_ = products
	/*
		tx, err := s.database.Begin()
		if err != nil {
			return err
		}
		defer func() { _ = tx.Rollback() }()

		if err = s.createTasks(products); err != nil {
			return err
		}
		if err = s.addToBuffer(products); err != nil {
			return err
		}
		if err = s.addToDatabase(tx, products); err != nil {
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
	*/
	return nil
}

func (s *Service) createTasks(products []*Product) error {
	// TODO что она будет делать ?
	return nil
}

func (s *Service) addToBuffer(products []*Product) error {
	// TODO узнать, явзяется ли номер карточки уникальным !?
	// TODO надо ли проверять наличие ?
	//for _, v := range products {
	//	if _, isExist := s.products[v.NumberOfCard]; isExist {
	//		return fmt.Errorf("error: product %d exist", v.NumberOfCard)
	//	}
	//}
	for _, v := range products {
		s.products[v.NumberOfCard] = v
	}
	return nil
}

func (s *Service) createTransaction() (*sql.Tx, error) {
	return nil, nil
}

func (s *Service) addToDatabase(tx *sql.Tx, products []*Product) error {
	return nil
}

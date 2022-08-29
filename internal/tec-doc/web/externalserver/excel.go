package externalserver

import (
	"errors"
	exl "github.com/xuri/excelize/v2"
	"io"
	"strconv"
	"tec-doc/pkg/model"
)

func (e *externalHttpServer) loadFromExcel(bodyData io.Reader) (products []model.Product, err error) {
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

func (e *externalHttpServer) parseExcelRow(p *model.Product, row []string) (err error) {
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
	return nil
}

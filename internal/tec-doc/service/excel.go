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
	f.SetSheetName(f.GetSheetName(0), nameSheet)
	// Set values
	_ = f.SetCellValue(nameSheet, "A1", "Номер карточки")
	_ = f.SetCellValue(nameSheet, "B1", "Артикул поставщика (уникальный артикул)")
	_ = f.SetCellValue(nameSheet, "C1", "Артикул производителя (артикул tec-doc)")
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

package service

import (
	exl "github.com/xuri/excelize/v2"
)

func (s *Service) ExcelTemplateForCLient(nameSheet string) ([]byte, error) {
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

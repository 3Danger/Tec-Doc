package service

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	exl "github.com/xuri/excelize/v2"
	"testing"
)

func TestService_ExcelTemplateForClient(t *testing.T) {
	nameSheet := "Products"
	s := Service{}
	xlsBytes, err := s.ExcelTemplateForClient()
	if !assert.NoError(t, err) {
		return
	}
	if !assert.NotEmpty(t, xlsBytes) {
		return
	}
	xlsReader := bytes.NewReader(xlsBytes)
	xlsFile, err := exl.OpenReader(xlsReader)
	defer func() {
		_ = xlsFile.Close()
	}()
	assert.Equal(t, 1, xlsFile.SheetCount, "sheet list too many")
	if !assert.NotEqual(t, -1, xlsFile.GetSheetIndex(nameSheet), "sheet \""+nameSheet+"\" not found") {
		return
	}
	if !assert.Equal(t, nameSheet, xlsFile.GetSheetName(xlsFile.GetActiveSheetIndex()), nameSheet+" sheet is not active") {
		xlsFile.SetActiveSheet(xlsFile.GetSheetIndex(nameSheet))
	}

	cellChecker := func(axis, compare string) {
		value, err := xlsFile.GetCellValue(nameSheet, axis)
		assert.NoError(t, err)
		assert.Equal(t, compare, value)
	}
	cellChecker("A1", "Номер карточки")
	cellChecker("B1", "Артикул поставщика (уникальный артикул)")
	cellChecker("C1", "Артикул производителя")
	cellChecker("D1", "Бренд")
	cellChecker("E1", "SKU")
	cellChecker("F1", "Категория товара")
	cellChecker("G1", "Цена товара")
}

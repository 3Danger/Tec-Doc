package service

import (
	exl "github.com/xuri/excelize/v2"
	"strconv"
)

type DummyXLSX struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

func ToExcel(filepath string, nameSheet string, data []DummyXLSX) (err error) {
	if data == nil || len(data) == 0 {
		return nil
	}

	f := exl.NewFile()
	defer func() { _ = f.Close() }()

	if nameSheet == "" {
		nameSheet = "Layer1"
	}

	index := f.NewSheet(nameSheet)
	_ = f.SetCellValue(nameSheet, "A1", "Id")
	_ = f.SetCellValue(nameSheet, "B1", "Name")
	_ = f.SetCellValue(nameSheet, "C1", "Price")
	for i, v := range data {
		_ = f.SetCellValue(nameSheet, "A"+strconv.Itoa(i+2), v.ID)
		_ = f.SetCellValue(nameSheet, "B"+strconv.Itoa(i+2), v.Name)
		_ = f.SetCellFloat(nameSheet, "C"+strconv.Itoa(i+2), v.Price, 20, 64)
	}
	f.SetActiveSheet(index)
	return f.SaveAs(filepath + ".xlsx")
}

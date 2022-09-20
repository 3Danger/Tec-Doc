package service

import (
	"bytes"
	"context"
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

func setExcelHeader(stream *exl.StreamWriter, style int, headers ...string) (err error) {
	if len(headers) == 0 {
		return errors.New("haven't headers")
	}
	// Set length of cell
	for i, w := range headers {
		if err = stream.SetColWidth(i+1, i+1, float64(len(w))*0.7); err != nil {
			return err
		}
	}

	// Set values
	var _headers = make([]interface{}, 0, len(headers))
	for i := range headers {
		_headers = append(_headers, headers[i])
	}
	return stream.SetRow("A1", _headers, exl.RowOpts{20, false, style})
}

func (s *Service) ExcelTemplateForClient() ([]byte, error) {
	var (
		err  error
		file = exl.NewFile()
	)
	defer func() { _ = file.Close() }()

	file.SetSheetName(file.GetSheetName(0), "Продукты")
	var stream *exl.StreamWriter
	if stream, err = file.NewStreamWriter(file.GetSheetName(0)); err != nil {
		return nil, err
	}

	var headers = []string{
		"Бренд",
		"Артикул поставщика (уникальный артикул)",
		"Артикул производителя (артикул tec-doc)",
		"Цена товара",
		"Штрих-код",
		"Комплектация",
	}
	if err = setExcelHeader(stream, 0, headers...); err != nil {
		return nil, err
	}

	var buffer *bytes.Buffer
	if buffer, err = file.WriteToBuffer(); err != nil {
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
	p.Amount = 1
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

	uploaderId, err := s.database.CreateTask(ctx, tx, ctx.GetString("X-Supplier-Id"), supplierID, userID, int64(len(products)), ctx.ClientIP(), time.Now().UTC())
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

	var file = exl.NewFile()
	defer func() { _ = file.Close() }()
	file.SetSheetName(file.GetSheetName(0), "Детализация продуктов")
	var stream *exl.StreamWriter
	if stream, err = file.NewStreamWriter(file.GetSheetName(0)); err != nil {
		return nil, err
	}

	var headers = []string{
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
	}

	if err = setExcelHeader(stream, 0, headers...); err != nil {
		return nil, err
	}

	var style int
	if style, err = file.NewStyle(styleExcel); err != nil {
		return nil, err
	}
	for i, p := range productsEnriched {
		ch := s.tecDocClient.ConvertToCharacteristics(&p)
		axis := fmt.Sprintf("A%d", i+2)
		err = stream.SetRow(axis, []interface{}{
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
	if err = stream.Flush(); err != nil {
		return nil, err
	}
	var buffer *bytes.Buffer
	if buffer, err = file.WriteToBuffer(); err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func (s *Service) ExcelProductsHistoryWithStatus(ctx context.Context, uploadID, status int64) ([]byte, error) {
	var (
		products []model.Product
		err      error
	)
	//TODO ? узнать, необходимо ли ограничение
	if products, err = s.database.GetProductsHistoryWithStatus(ctx, nil, uploadID, status, 100000, 0); err != nil {
		return nil, err
	}

	var file = exl.NewFile()
	defer func() { _ = file.Close() }()
	file.SetSheetName(file.GetSheetName(0), "Продукты с ошибками")
	var stream *exl.StreamWriter
	if stream, err = file.NewStreamWriter(file.GetSheetName(0)); err != nil {
		return nil, err
	}

	var headers = []string{
		"Бренд",
		"Артикул поставщика (уникальный артикул)",
		"Артикул производителя (артикул tec-doc)",
		"Цена товара",
		"Штрих-код",
		"Комплектация",
		"Ошибка при обработки",
	}
	if err = setExcelHeader(stream, 0, headers...); err != nil {
		return nil, err
	}

	for i := range products {
		var axis string
		if axis, err = exl.CoordinatesToCellName(1, i+2); err != nil {
			return nil, err
		}
		if err = stream.SetRow(axis, []interface{}{
			products[i].Brand,
			products[i].ArticleSupplier,
			products[i].Article,
			products[i].Price,
			products[i].Barcode,
			products[i].Amount,
			products[i].ErrorResponse,
		}); err != nil {
			return nil, err
		}
	}
	if err = stream.Flush(); err != nil {
		return nil, err
	}
	var buf *bytes.Buffer
	if buf, err = file.WriteToBuffer(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

package postgres

import (
	"context"
	"errors"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"strings"
	"tec-doc/internal/tec-doc/model"
	"testing"
	"time"

	//"tec-doc/internal/tec-doc/mocks/mock_postgres"
	_ "github.com/pashagolub/pgxmock"
)

//getNameFields возвращает список с именами полей
//func getNameFields(value interface{}) []string {
//	valueType := reflect.TypeOf(value)
//	length := valueType.NumField()
//
//	var fields = make([]string, length)
//	for i := 0; i < length; i++ {
//		fields[i] = strings.ToLower(valueType.Field(i).Name)
//	}
//	return fields
//}

func TestStore_GetProductsHistory(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockPool.Close()

	type args struct {
		db       Pool
		ctx      context.Context
		tx       Transaction
		UploadID int64
		limit    int
		offset   int
	}
	type mockBehavior func(args args, products []model.Product, wantError bool)
	var namesFields = []string{
		"id", "upload_id", "article", "card_number", "provider_article", "manufacturer_article",
		"brand", "sku", "category", "price", "upload_date", "update_date", "status", "errorresponse",
	}
	var behavior = func(args args, products []model.Product, wantError bool) {
		var err error = nil
		if wantError {
			err = errors.New("some err")
		}
		rows := pgxmock.NewRows(namesFields)
		for _, p := range products {
			rows.AddRow(p.ID, p.UploadID, p.Article, p.CardNumber, p.ProviderArticle, p.ManufacturerArticle,
				p.Brand, p.SKU, p.Category, p.Price, p.UploadDate, p.UpdateDate, p.Status, "")
		}
		mockPool.ExpectQuery("SELECT "+strings.Join(namesFields, ", ")+
			" FROM products_history WHERE upload_id = \\$1 LIMIT \\$2 OFFSET \\$3").
			WithArgs(args.UploadID, args.limit, args.offset).WillReturnRows(rows).WillReturnError(err)
	}

	testCases := map[string]struct {
		input   args
		wantErr bool
		want    []model.Product
		mock    mockBehavior
	}{
		"OK pool executor": {
			input:   args{mockPool, context.Background(), nil, 2, 1, 0},
			wantErr: false,
			want: []model.Product{
				{
					1, 2, "3", 4, "5",
					"6", "7", "8", "9", 10,
					time.Now(), time.Now(), 11, "",
				},
			},
			mock: behavior,
		},
		"OK pool tx": {
			input:   args{nil, context.Background(), mockPool, 2, 1, 0},
			wantErr: false,
			want: []model.Product{
				{
					1, 2, "3", 4, "5",
					"6", "7", "8", "9", 10,
					time.Now(), time.Now(), 11, "",
				},
			},
			mock: behavior,
		},
		"empty rows executor": {
			input:   args{mockPool, context.Background(), nil, 2, 1, 0},
			wantErr: false,
			want:    nil,
			mock:    behavior,
		},
		"error pool executor": {
			input:   args{mockPool, context.Background(), nil, 2, 1, 0},
			wantErr: true,
			want:    nil,
			mock:    behavior,
		},
	}
	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			tt.mock(tt.input, tt.want, tt.wantErr)
			got, err := (&store{pool: tt.input.db}).GetProductsHistory(tt.input.ctx, tt.input.tx, tt.input.UploadID, tt.input.limit, tt.input.offset)
			assert.NoError(t, mockPool.ExpectationsWereMet())
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

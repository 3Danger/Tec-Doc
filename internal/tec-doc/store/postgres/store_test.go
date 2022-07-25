package postgres

import (
	"context"
	"errors"
	"tec-doc/internal/tec-doc/model"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strings"
	"tec-doc/internal/tec-doc/config"
)

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
					1, 2, 3, "4", "5",
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
					1, 2, 3, "4", "5",
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

func TestStore_GetSupplierTaskHistory(t *testing.T) {
	mockPool, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer mockPool.Close()

	mockStore := store{
		cfg:  &config.PostgresConfig{},
		pool: mockPool,
	}

	type args struct {
		ctx        context.Context
		tx         Transaction
		supplierID int64
		limit      int
		offset     int
	}
	type mockBehavior func(args args, tasks []model.Task)

	testCases := []struct {
		name    string
		input   args
		want    []model.Task
		mock    mockBehavior
		wantErr bool
	}{
		{
			name:  "OK pool executor",
			input: args{context.Background(), nil, int64(1), 0, 0},
			want:  []model.Task{{int64(1), int64(1), int64(1), time.Now(), time.Now(), "127.0.0.1", 0, 0, 0, 0}},
			mock: func(args args, tasks []model.Task) {
				t := tasks[0]
				rows := mockPool.NewRows([]string{"id", "supplier_id", "user_id", "IP",
					"upload_date", "update_date", "status", "products_processed", "products_failed", "products_total"}).AddRow(
					t.ID, t.SupplierID, t.UserID, t.IP, t.UploadDate, t.UpdateDate, t.Status, t.ProductsProcessed, t.ProductsFailed, t.ProductsTotal)
				mockPool.ExpectQuery(regexp.QuoteMeta(`SELECT id, supplier_id, user_id, IP, upload_date, update_date, status, products_processed, products_failed, products_total FROM tasks WHERE supplier_id = $1 ORDER BY upload_date LIMIT $2 OFFSET $3;`)).WithArgs(args.supplierID, args.limit, args.offset).WillReturnRows(rows)
			},
		},
		{
			name:  "OK tx executor",
			input: args{context.Background(), mockPool, int64(1), 0, 0},
			want:  []model.Task{{int64(1), int64(1), int64(1), time.Now(), time.Now(), "127.0.0.1", 0, 0, 0, 0}},
			mock: func(args args, tasks []model.Task) {
				t := tasks[0]
				rows := mockPool.NewRows([]string{"id", "supplier_id", "user_id", "IP",
					"upload_date", "update_date", "status", "products_processed", "products_failed", "products_total"}).AddRow(
					t.ID, t.SupplierID, t.UserID, t.IP, t.UploadDate, t.UpdateDate, t.Status, t.ProductsProcessed, t.ProductsFailed, t.ProductsTotal)
				mockPool.ExpectQuery(regexp.QuoteMeta(`SELECT id, supplier_id, user_id, IP, upload_date, update_date, status, products_processed, products_failed, products_total FROM tasks WHERE supplier_id = $1 ORDER BY upload_date LIMIT $2 OFFSET $3;`)).WithArgs(args.supplierID, args.limit, args.offset).WillReturnRows(rows)
			},
		},
		{
			name:  "Error wrong sql query",
			input: args{context.Background(), mockPool, int64(1), 0, 0},
			want:  nil,
			mock: func(args args, tasks []model.Task) {
				mockPool.ExpectQuery(regexp.QuoteMeta(`id, supplier_id, user_id, IP, upload_date, update_date, status, products_processed, products_failed, products_total FROM tasks WHERE supplier_id = $1 ORDER BY upload_date LIMIT $2 OFFSET $3;`)).WithArgs(args.supplierID, args.limit, args.offset).WillReturnError(pgxmock.ErrCancelled)
			},
			wantErr: true,
		},
		{
			name:  "Error invalid arguments order when rows.Scan()",
			input: args{context.Background(), mockPool, int64(1), 0, 0},
			want:  nil,
			mock: func(args args, tasks []model.Task) {
				t := model.Task{int64(1), int64(1), int64(1), time.Now(), time.Now(), "127.0.0.1", 0, 0, 0, 0}
				rows := mockPool.NewRows([]string{"id", "supplier_id", "user_id", "IP",
					"upload_date", "update_date", "status", "products_processed", "products_failed", "products_total"}).AddRow(
					t.ID, t.SupplierID, t.IP, t.UserID, t.UploadDate, t.UpdateDate, t.Status, t.ProductsProcessed, t.ProductsFailed, t.ProductsTotal)
				mockPool.ExpectQuery(regexp.QuoteMeta(`SELECT id, supplier_id, user_id, IP, upload_date, update_date, status, products_processed, products_failed, products_total FROM tasks WHERE supplier_id = $1 ORDER BY upload_date LIMIT $2 OFFSET $3;`)).WithArgs(args.supplierID, args.limit, args.offset).WillReturnRows(rows)
			},
			wantErr: true,
		},
		{
			name:  "Error invalid arguments number when rows.Scan()",
			input: args{context.Background(), nil, int64(1), 0, 0},
			want:  []model.Task{{int64(1), int64(1), int64(1), time.Now(), time.Now(), "127.0.0.1", 0, 0, 0, 0}},
			mock: func(args args, tasks []model.Task) {
				rows := mockPool.NewRows([]string{"id", "supplier_id", "user_id", "IP",
					"upload_date", "update_date", "status", "products_processed", "products_failed", "products_total", "err"}).AddRow(
					tasks[0].ID, tasks[0].SupplierID, tasks[0].UserID, tasks[0].IP, tasks[0].UploadDate, tasks[0].UpdateDate, tasks[0].Status, tasks[0].ProductsProcessed, tasks[0].ProductsFailed, tasks[0].ProductsTotal, "err")
				mockPool.ExpectQuery(regexp.QuoteMeta(`SELECT id, supplier_id, user_id, IP, upload_date, update_date, status, products_processed, products_failed, products_total FROM tasks WHERE supplier_id = $1 ORDER BY upload_date LIMIT $2 OFFSET $3;`)).WithArgs(args.supplierID, args.limit, args.offset).WillReturnRows(rows)
			},
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input, tt.want)

			got, err := mockStore.GetSupplierTaskHistory(tt.input.ctx, tt.input.tx, tt.input.supplierID, tt.input.limit, tt.input.offset)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mockPool.ExpectationsWereMet())
		})
	}
}

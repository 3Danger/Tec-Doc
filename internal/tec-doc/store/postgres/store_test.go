package postgres

import (
	"context"
	"fmt"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"regexp"
	"tec-doc/internal/tec-doc/config"
	"tec-doc/internal/tec-doc/model"
	"testing"
	"time"
)

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
	type mockBehavoir func(args args, tasks []model.Task)

	testCases := []struct {
		name    string
		input   args
		want    []model.Task
		mock    mockBehavoir
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
				fmt.Println(err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mockPool.ExpectationsWereMet())
		})
	}
}

package service

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"tec-doc/internal/tec-doc/config"
	mock_service "tec-doc/internal/tec-doc/mock"
	"tec-doc/internal/tec-doc/model"
	"tec-doc/internal/tec-doc/store/postgres"
	"testing"
	"time"
)

func TestService_GetSupplierTaskHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_service.NewMockStore(ctrl)
	service := &Service{
		conf:     &config.Config{},
		log:      &zerolog.Logger{},
		database: mockStore,
	}

	type args struct {
		ctx        context.Context
		tx         postgres.Transaction
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
			name:  "OK",
			input: args{context.Background(), nil, int64(1), 0, 0},
			want:  []model.Task{{int64(1), int64(1), int64(1), time.Now(), time.Now(), "127.0.0.1", 0, 0, 0, 0}},
			mock: func(args args, tasks []model.Task) {
				mockStore.EXPECT().GetSupplierTaskHistory(args.ctx, args.tx, args.supplierID, args.limit, args.offset).Return(tasks, nil).Times(1)
			},
		},
		{
			name:  "Error invalid limit",
			input: args{context.Background(), nil, int64(1), -1, 0},
			want:  nil,
			mock: func(args args, tasks []model.Task) {
				mockStore.EXPECT().GetSupplierTaskHistory(args.ctx, args.tx, args.supplierID, args.limit, args.offset).Return(nil, nil).Times(1)
			},
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input, tt.want)

			got, err := service.GetSupplierTaskHistory(tt.input.ctx, tt.input.tx, tt.input.supplierID, tt.input.limit, tt.input.offset)
			if tt.wantErr {
				assert.Error(t, err)
				fmt.Println(err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

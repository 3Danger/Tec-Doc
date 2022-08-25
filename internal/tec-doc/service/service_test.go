package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"tec-doc/internal/tec-doc/config"
	"tec-doc/internal/tec-doc/mocks/mock_externalserver"
	"tec-doc/internal/tec-doc/mocks/mock_internalserver"
	"tec-doc/internal/tec-doc/mocks/mock_postgres"
	"tec-doc/internal/tec-doc/store/postgres"
	"tec-doc/pkg/clients/model"
	"testing"
	"time"
)

func TestService_GetSupplierTaskHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStore := mock_postgres.NewMockStore(ctrl)
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

func TestService_Start(t *testing.T) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	funcSleep := func() func() error {
		return func() error { time.Sleep(time.Second * 10); return nil }
	}
	svc := &Service{
		log: new(zerolog.Logger),
	}

	testCases := map[string]struct {
		err                    error
		contextDone            bool
		internalServerBehavior func(server *mock_internalserver.MockServer)
		externalServerBehavior func(server *mock_externalserver.MockServer)
	}{
		"success": {
			internalServerBehavior: func(server *mock_internalserver.MockServer) {
				server.EXPECT().Start().AnyTimes().Return(nil)
			},
			externalServerBehavior: func(server *mock_externalserver.MockServer) {
				server.EXPECT().Start().AnyTimes().Return(nil)
			},
		},
		"success context done": {
			contextDone: true,
			internalServerBehavior: func(server *mock_internalserver.MockServer) {
				server.EXPECT().Start().AnyTimes().DoAndReturn(funcSleep())
			},
			externalServerBehavior: func(server *mock_externalserver.MockServer) {
				server.EXPECT().Start().AnyTimes().DoAndReturn(funcSleep())
			},
		},
		"error internal": {
			err: errors.New("test error of internal-server"),
			internalServerBehavior: func(server *mock_internalserver.MockServer) {
				server.EXPECT().Start().AnyTimes().Return(errors.New("test error of internal-server"))
			},
			externalServerBehavior: func(server *mock_externalserver.MockServer) {
				server.EXPECT().Start().AnyTimes().Do(funcSleep())
			},
		},
		"error external": {
			err: errors.New("test error of external-server"),
			internalServerBehavior: func(server *mock_internalserver.MockServer) {
				server.EXPECT().Start().AnyTimes().Do(funcSleep())
			},
			externalServerBehavior: func(server *mock_externalserver.MockServer) {
				server.EXPECT().Start().Return(errors.New("test error of external-server"))
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			ctx, cancelFunc := context.WithCancel(context.Background())
			defer cancelFunc()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			internalServer := mock_internalserver.NewMockServer(ctrl)
			externalServer := mock_externalserver.NewMockServer(ctrl)
			testCase.internalServerBehavior(internalServer)
			testCase.externalServerBehavior(externalServer)
			svc.internalServer = internalServer
			svc.externalServer = externalServer
			if testCase.contextDone {
				cancelFunc()
			}
			assert.Equal(t, testCase.err, svc.Start(ctx))
		})
	}
}

func TestService_Stop(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockInternalServer := mock_internalserver.NewMockServer(ctrl)
	mockExternalServer := mock_externalserver.NewMockServer(ctrl)
	zerolog.SetGlobalLevel(zerolog.Disabled)
	service := Service{
		log:            new(zerolog.Logger),
		externalServer: mockInternalServer,
		internalServer: mockExternalServer,
	}

	//Just stop
	mockInternalServer.EXPECT().Stop()
	mockExternalServer.EXPECT().Stop()
	service.Stop()

	//Stop with return error
	mockInternalServer.EXPECT().Stop().Return(errors.New("some error"))
	mockExternalServer.EXPECT().Stop().Return(errors.New("some error"))
	service.Stop()
}

package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"tec-doc/internal/tec-doc/mocks/mock_externalserver"
	"tec-doc/internal/tec-doc/mocks/mock_internalserver"
	"testing"
	"time"
)

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

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
	doAndReturnError := func(second time.Duration, errMsg string) func() error {
		var err error
		if errMsg != "" {
			err = errors.New(errMsg)
		}
		return func() error { time.Sleep(time.Second * second); return err }
	}
	svc := &Service{
		log: new(zerolog.Logger),
	}

	testCases := map[string]struct {
		err                    error
		internalServerBehavior func(server *mock_internalserver.MockServer)
		externalServerBehavior func(server *mock_externalserver.MockServer)
	}{
		"success": {
			err: nil,
			internalServerBehavior: func(server *mock_internalserver.MockServer) {
				server.EXPECT().Start().DoAndReturn(doAndReturnError(0, ""))
			},
			externalServerBehavior: func(server *mock_externalserver.MockServer) {
				server.EXPECT().Start().DoAndReturn(doAndReturnError(10, ""))
			},
		},
		"error internal": {
			err: errors.New("test error of internal-server"),
			internalServerBehavior: func(server *mock_internalserver.MockServer) {
				server.EXPECT().Start().DoAndReturn(doAndReturnError(0, "test error of internal-server"))
			},
			externalServerBehavior: func(server *mock_externalserver.MockServer) {
				server.EXPECT().Start().Do(doAndReturnError(10, ""))
			},
		},
		"error external": {
			err: errors.New("test error of external-server"),
			internalServerBehavior: func(server *mock_internalserver.MockServer) {
				server.EXPECT().Start().Do(doAndReturnError(10, ""))
			},
			externalServerBehavior: func(server *mock_externalserver.MockServer) {
				server.EXPECT().Start().DoAndReturn(doAndReturnError(0, "test error of external-server"))
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			internalServer := mock_internalserver.NewMockServer(ctrl)
			externalServer := mock_externalserver.NewMockServer(ctrl)
			testCase.internalServerBehavior(internalServer)
			testCase.externalServerBehavior(externalServer)
			svc.internalServer = internalServer
			svc.externalServer = externalServer
			err := svc.Start(context.TODO())
			assert.Equal(t, testCase.err, err)
		})
	}
}

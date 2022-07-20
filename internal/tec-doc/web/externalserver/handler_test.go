package externalserver

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	mock_externalserver "tec-doc/internal/tec-doc/mock"
	"tec-doc/internal/tec-doc/model"
	"tec-doc/internal/tec-doc/store/postgres"
	"tec-doc/internal/tec-doc/web/externalserver/middleware"
	"testing"
	"time"
)

func TestHandler_GetSupplierTaskHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mock_externalserver.NewMockService(ctrl)
	router := gin.Default()
	server := externalHttpServer{
		router:  router,
		service: mockService,
		logger:  &zerolog.Logger{},
		server: http.Server{
			Addr:    "/",
			Handler: router,
		},
	}
	server.router.GET("/task_history_test", middleware.Authorize, server.GetSupplierTaskHistory)

	type args struct {
		ctx              context.Context
		tx               postgres.Transaction
		userID           string
		supplierID       string
		limit            string
		offset           string
		userIDHeader     string
		supplierIDHeader string
		limitQuery       string
		offsetQuery      string
	}

	type mockBehavoir func(args args, tasks []model.Task)

	testCases := []struct {
		name                 string
		input                args
		mock                 mockBehavoir
		wantOkBody           []model.Task
		wantErrBody          string
		expectedStatusCode   int
		expectedResponseBody func(want interface{}) string
	}{
		{
			name:       "OK",
			input:      args{context.Background(), nil, "1", "1", "0", "0", "X-User-Id", "X-Supplier-Id", "limit", "offset"},
			wantOkBody: []model.Task{{int64(1), int64(1), int64(1), time.Now(), time.Now(), "127.0.0.1", 0, 0, 0, 0}},
			mock: func(args args, tasks []model.Task) {
				supplierID, _ := strconv.ParseInt(args.supplierID, 10, 64)
				limit, _ := strconv.Atoi(args.limit)
				offset, _ := strconv.Atoi(args.offset)
				mockService.EXPECT().GetSupplierTaskHistory(gomock.Any(), args.tx, supplierID, limit, offset).Return(tasks, nil).Times(1)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: func(want interface{}) string {
				js, _ := json.Marshal(want)
				return string(js)
			},
		},
		{
			name:        "Error invalid userID header",
			input:       args{context.Background(), nil, "1", "1", "0", "0", "Error", "X-Supplier-Id", "limit", "offset"},
			wantErrBody: `{"error":"invalid user_id"}{"error":"can't get user or supplier id from context"}`,
			mock: func(args args, tasks []model.Task) {
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponseBody: func(want interface{}) string {
				return want.(string)
			},
		},
		{
			name:        "Error invalid supplierID header",
			input:       args{context.Background(), nil, "1", "1", "0", "0", "X-User-Id", "Error", "limit", "offset"},
			wantErrBody: `{"error":"invalid supplier_id"}{"error":"can't get user or supplier id from context"}`,
			mock: func(args args, tasks []model.Task) {
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponseBody: func(want interface{}) string {
				return want.(string)
			},
		},
		{
			name:        "Error invalid limit query",
			input:       args{context.Background(), nil, "1", "1", "0", "0", "X-User-Id", "X-Supplier-Id", "Error", "offset"},
			wantErrBody: `{"can't get limit":"strconv.Atoi: parsing \"\": invalid syntax"}`,
			mock: func(args args, tasks []model.Task) {
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: func(want interface{}) string {
				return want.(string)
			},
		},
		{
			name:        "Error invalid offset query",
			input:       args{context.Background(), nil, "1", "1", "0", "0", "X-User-Id", "X-Supplier-Id", "limit", "Error"},
			wantErrBody: `{"can't get offset":"strconv.Atoi: parsing \"\": invalid syntax"}`,
			mock: func(args args, tasks []model.Task) {
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponseBody: func(want interface{}) string {
				return want.(string)
			},
		},
		{
			name:        "Error invalid limit value",
			input:       args{context.Background(), nil, "1", "1", "-1", "0", "X-User-Id", "X-Supplier-Id", "limit", "offset"},
			wantErrBody: `{"error":"{can't get task history}"}`,
			mock: func(args args, tasks []model.Task) {
				supplierID, _ := strconv.ParseInt(args.supplierID, 10, 64)
				limit, _ := strconv.Atoi(args.limit)
				offset, _ := strconv.Atoi(args.offset)
				mockService.EXPECT().GetSupplierTaskHistory(gomock.Any(), args.tx, supplierID, limit, offset).Return(nil, errors.New("{can't get task history}")).Times(1)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponseBody: func(want interface{}) string {
				return want.(string)
			},
		},
		{
			name:        "Error invalid offset value",
			input:       args{context.Background(), nil, "1", "1", "0", "-1", "X-User-Id", "X-Supplier-Id", "limit", "offset"},
			wantErrBody: `{"error":"{can't get task history}"}`,
			mock: func(args args, tasks []model.Task) {
				supplierID, _ := strconv.ParseInt(args.supplierID, 10, 64)
				limit, _ := strconv.Atoi(args.limit)
				offset, _ := strconv.Atoi(args.offset)
				mockService.EXPECT().GetSupplierTaskHistory(gomock.Any(), args.tx, supplierID, limit, offset).Return(nil, errors.New("{can't get task history}")).Times(1)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponseBody: func(want interface{}) string {
				return want.(string)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input, tt.wantOkBody)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", "/task_history_test", nil)
			req.Header = http.Header{tt.input.userIDHeader: {tt.input.userID}, tt.input.supplierIDHeader: {tt.input.supplierID}}
			req.URL.RawQuery = url.Values{tt.input.limitQuery: {tt.input.limit}, tt.input.offsetQuery: {tt.input.offset}}.Encode()

			server.router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			if tt.wantErrBody == "" {
				assert.Equal(t, w.Body.String(), tt.expectedResponseBody(tt.wantOkBody))
			} else {
				assert.Equal(t, w.Body.String(), tt.expectedResponseBody(tt.wantErrBody))
			}
		})
	}
}

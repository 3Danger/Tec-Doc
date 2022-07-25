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
	mockExternalServer "tec-doc/internal/tec-doc/mocks/mock_externalserver"
	"tec-doc/internal/tec-doc/model"
	"tec-doc/internal/tec-doc/store/postgres"
	"tec-doc/internal/tec-doc/web/externalserver/middleware"
	"testing"
	"time"
)

func TestExternalHttpServer_ExcelTemplate(t *testing.T) {
	type mockBehavior func(service *mockExternalServer.MockService)
	wantErr := errors.New("some error")
	wantErrBytes, _ := json.Marshal(map[string]string{"error": wantErr.Error()})
	testCase := map[string]struct {
		contentType string
		behavior    mockBehavior
		want        []byte
		wantStatus  int
	}{
		"success test": {
			contentType: "application/vnd.ms-excel",
			behavior: func(service *mockExternalServer.MockService) {
				service.EXPECT().ExcelTemplateForClient().Return([]byte("some valid excel data"), nil)
			},
			want:       []byte("some valid excel data"),
			wantStatus: http.StatusOK,
		},
		"error test": {
			contentType: "application/json",
			behavior: func(service *mockExternalServer.MockService) {
				service.EXPECT().ExcelTemplateForClient().Return(nil, wantErr)
			},
			want:       wantErrBytes,
			wantStatus: http.StatusInternalServerError,
		},
	}
	gin.SetMode(gin.TestMode)
	zerolog.SetGlobalLevel(zerolog.Disabled)
	for name, tc := range testCase {
		t.Run(name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockService := mockExternalServer.NewMockService(ctrl)
			tc.behavior(mockService)
			server := externalHttpServer{
				router:  gin.New(),
				service: mockService,
				logger:  new(zerolog.Logger),
			}
			server.router.GET("/excel_template", server.ExcelTemplate)

			req := httptest.NewRequest(http.MethodGet, "/excel_template", nil)
			w := httptest.NewRecorder()

			server.router.ServeHTTP(w, req)
			contentType := w.Header().Get("Content-Type")
			assert.Equal(t, tc.wantStatus, w.Code)
			assert.Contains(t, contentType, tc.contentType)
			assert.Equal(t, tc.want, w.Body.Bytes())
		})
	}
}

func TestExternalHttpServer_GetSupplierTaskHistory(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockService := mockExternalServer.NewMockService(ctrl)

	type args struct {
		ctx              context.Context
		tx               postgres.Transaction
		userID           int64
		supplierID       int64
		limit            int
		offset           int
		userIDHeader     string
		supplierIDHeader string
		limitQuery       string
		offsetQuery      string
	}

	type mockBehavior func(args args, tasks []model.Task)

	testCases := []struct {
		name                 string
		input                args
		tasks                []model.Task
		mock                 mockBehavior
		expectedStatusCode   int
		wantErrBody          string
		expectedResponseBody func(want interface{}) string
	}{
		{
			name:  "OK",
			input: args{context.Background(), nil, int64(1), int64(1), 0, 0, "X-User-Id", "X-Supplier-Id", "limit", "offset"},
			tasks: []model.Task{{int64(1), int64(1), int64(1), time.Now(), time.Now(), "127.0.0.1", 0, 0, 0, 0}},
			mock: func(args args, tasks []model.Task) {
				mockService.EXPECT().GetSupplierTaskHistory(gomock.Any(), args.tx, args.supplierID, args.limit, args.offset).Return(tasks, nil).Times(1)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponseBody: func(want interface{}) string {
				js, _ := json.Marshal(want)
				return string(js)
			},
		},
		{
			name:        "Error invalid userID header",
			input:       args{context.Background(), nil, int64(1), int64(1), 0, 0, "Error", "X-Supplier-Id", "limit", "offset"},
			wantErrBody: `{"error":"can't get user_id from context"}`,
			mock: func(args args, tasks []model.Task) {
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponseBody: func(want interface{}) string {
				return want.(string)
			},
		},
		{
			name:        "Error invalid supplierID header",
			input:       args{context.Background(), nil, int64(1), int64(1), 0, 0, "X-User-Id", "Error", "limit", "offset"},
			wantErrBody: `{"error":"can't get supplier_id from context"}`,
			mock: func(args args, tasks []model.Task) {
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponseBody: func(want interface{}) string {
				return want.(string)
			},
		},
		{
			name:        "Error invalid limit query",
			input:       args{context.Background(), nil, int64(1), int64(1), 0, 0, "X-User-Id", "X-Supplier-Id", "Error", "offset"},
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
			input:       args{context.Background(), nil, int64(1), int64(1), 0, 0, "X-User-Id", "X-Supplier-Id", "limit", "Error"},
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
			input:       args{context.Background(), nil, int64(1), int64(1), -1, 0, "X-User-Id", "X-Supplier-Id", "limit", "offset"},
			wantErrBody: `{"error":"{can't get task history}"}`,
			mock: func(args args, tasks []model.Task) {
				mockService.EXPECT().GetSupplierTaskHistory(gomock.Any(), args.tx, args.supplierID, args.limit, args.offset).Return(nil, errors.New("{can't get task history}")).Times(1)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponseBody: func(want interface{}) string {
				return want.(string)
			},
		},
		{
			name:        "Error invalid offset value",
			input:       args{context.Background(), nil, int64(1), int64(1), 1, -1, "X-User-Id", "X-Supplier-Id", "limit", "offset"},
			wantErrBody: `{"error":"{can't get task history}"}`,
			mock: func(args args, tasks []model.Task) {
				mockService.EXPECT().GetSupplierTaskHistory(gomock.Any(), args.tx, args.supplierID, args.limit, args.offset).Return(nil, errors.New("{can't get task history}")).Times(1)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponseBody: func(want interface{}) string {
				return want.(string)
			},
		},
	}

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

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input, tt.tasks)

			w := httptest.NewRecorder()

			req := httptest.NewRequest("GET", "/task_history_test", nil)
			req.Header = http.Header{tt.input.userIDHeader: {strconv.FormatInt(tt.input.userID, 10)}, tt.input.supplierIDHeader: {strconv.FormatInt(tt.input.supplierID, 10)}}
			req.URL.RawQuery = url.Values{tt.input.limitQuery: {strconv.Itoa(tt.input.limit)}, tt.input.offsetQuery: {strconv.Itoa(tt.input.offset)}}.Encode()

			server.router.ServeHTTP(w, req)

			assert.Equal(t, w.Code, tt.expectedStatusCode)
			if tt.wantErrBody == "" {
				assert.Equal(t, tt.expectedResponseBody(tt.tasks), w.Body.String())
			} else {
				assert.Equal(t, tt.expectedResponseBody(tt.wantErrBody), w.Body.String())
			}
		})
	}
}

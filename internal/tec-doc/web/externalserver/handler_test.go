package externalserver

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	mockExternalServer "tec-doc/internal/tec-doc/mocks/mock_externalserver"
	"tec-doc/internal/tec-doc/model"
	"tec-doc/internal/tec-doc/web/externalserver/middleware"
	"testing"
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
	type mockBehavior func(svc *mockExternalServer.MockService, limit string, err error)
	var testCases = map[string]struct {
		userId     string
		supplierId string
		limit      string
		offset     string
		behavior   mockBehavior
		response   map[string]string
	}{
		"success test 1": {
			userId:     "9",
			supplierId: "10",
			limit:      "10",
			offset:     "10",
			behavior: func(svc *mockExternalServer.MockService, limit string, err error) {
				lim, _ := strconv.Atoi(limit)
				svc.EXPECT().
					GetSupplierTaskHistory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(make([]model.Task, lim), nil)
			},
			response: make(map[string]string),
		},
		"user_id error": {
			userId:     "f",
			supplierId: "10",
			limit:      "10",
			offset:     "10",
			behavior:   func(svc *mockExternalServer.MockService, limit string, err error) {},
			response:   map[string]string{"error": "can't get user_id from context"},
		},
		"supplier_id error": {
			userId:     "10",
			supplierId: "",
			limit:      "10",
			offset:     "10",
			behavior:   func(svc *mockExternalServer.MockService, limit string, err error) {},
			response:   map[string]string{"error": "can't get supplier_id from context"},
		},
		"limit error": {
			userId:     "9",
			supplierId: "10",
			limit:      "INVALID",
			offset:     "10",
			behavior:   func(svc *mockExternalServer.MockService, limit string, err error) {},
			response:   map[string]string{"error": "can't get limit"},
		},
		"offset error": {
			userId:     "9",
			supplierId: "10",
			limit:      "10",
			offset:     "INVALID",
			behavior:   func(svc *mockExternalServer.MockService, limit string, err error) {},
			response:   map[string]string{"error": "can't get offset"},
		},
		"GetSupplierTaskHistory error": {
			userId:     "9",
			supplierId: "10",
			limit:      "10",
			offset:     "10",
			behavior: func(svc *mockExternalServer.MockService, limit string, err error) {
				var tasksNil []model.Task = nil
				svc.EXPECT().
					GetSupplierTaskHistory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(tasksNil, errors.New("some error"))
			},
			response: map[string]string{"error": "some error"},
		},
	}

	gin.SetMode(gin.TestMode)
	zerolog.SetGlobalLevel(zerolog.Disabled)

	c := gomock.NewController(t)
	mockService := mockExternalServer.NewMockService(c)
	server := &externalHttpServer{
		service: mockService,
		logger:  &zerolog.Logger{},
	}

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			testCase.behavior(mockService, testCase.limit, errors.New(testCase.response["error"]))

			r := gin.New()
			r.Use(middleware.Authorize)
			r.GET("/tasks_history", server.GetSupplierTaskHistory)

			req := httptest.NewRequest(http.MethodGet, "/tasks_history", nil)
			q := req.URL.Query()
			q.Add("limit", testCase.limit)
			q.Add("offset", testCase.offset)
			req.URL.RawQuery = q.Encode()
			req.Header = http.Header{
				"X-User-Id":     {testCase.userId},
				"X-Supplier-Id": {testCase.supplierId},
			}

			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			var gotMap map[string]string
			_ = json.NewDecoder(w.Body).Decode(&gotMap)
			assert.Equal(t, testCase.response["error"], gotMap["error"])
		})
	}
}

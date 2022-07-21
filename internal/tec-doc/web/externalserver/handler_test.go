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
	"tec-doc/internal/tec-doc/model"
	"tec-doc/internal/tec-doc/web/externalserver/middleware"
	mockExternalServer "tec-doc/internal/tec-doc/web/externalserver/mocks"
	"testing"
)

func TestExternalHttpServer_GetSupplierTaskHistory(t *testing.T) {
	type mockBehavior func(svc *mockExternalServer.MockService, limit string, err error)
	simpleBehavior := func(svc *mockExternalServer.MockService, limit string, err error) {
		lim, _ := strconv.Atoi(limit)
		svc.EXPECT().
			GetSupplierTaskHistory(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(make([]model.Task, lim), nil)
	}
	var testCases = map[string]struct {
		userId     string
		supplierId string
		limit      string
		offset     string
		behavior   mockBehavior
		response   gin.H
	}{
		"Test1": {
			userId:     "9",
			supplierId: "10",
			limit:      "10",
			offset:     "10",
			behavior:   simpleBehavior,
			response:   gin.H{},
		},
		"Test2": {
			userId:     "9",
			supplierId: "10",
			limit:      "NOT VALID",
			offset:     "10",
			response:   gin.H{"error": "can't get limit"},
		},
		"Test3": {
			userId:     "9",
			supplierId: "10",
			limit:      "10",
			offset:     "NOT VALID",
			response:   gin.H{"error": "can't get offset"},
		},
		"Test4": {
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
			response: gin.H{"error": "some error"},
		},
	}

	gin.SetMode(gin.TestMode)
	c := gomock.NewController(t)
	defer c.Finish()
	service := mockExternalServer.NewMockService(c)

	fakeLogger := new(zerolog.Logger)
	fakeLogger.Level(zerolog.Disabled)
	server := New("8080", service, fakeLogger)

	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			r := gin.New()
			r.Use(middleware.Authorize)
			if testCase.behavior != nil {
				msg := ""
				if errMsg, ok := testCase.response["error"]; ok {
					msg = errMsg.(string)
				}
				testCase.behavior(service, testCase.limit, errors.New(msg))
			}
			r.GET("/tasks_history", server.GetSupplierTaskHistory)

			w := httptest.NewRecorder()

			req := httptest.NewRequest(http.MethodGet, "/tasks_history", nil)
			q := req.URL.Query()
			q.Add("limit", testCase.limit)
			q.Add("offset", testCase.offset)
			req.URL.RawQuery = q.Encode()
			req.Header = http.Header{
				"X-User-Id":     {testCase.userId},
				"X-Supplier-Id": {testCase.supplierId},
			}

			r.ServeHTTP(w, req)
			var writerGinMap gin.H
			_ = json.NewDecoder(w.Body).Decode(&writerGinMap)
			assert.Equal(t, testCase.response["error"], writerGinMap["error"])
		})
	}
}

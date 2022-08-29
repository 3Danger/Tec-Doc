// Code generated by MockGen. DO NOT EDIT.
// Source: internal/tec-doc/web/externalserver/server.go

// Package mock_externalserver is a generated GoMock package.
package mock_externalserver

import (
	context "context"
	io "io"
	reflect "reflect"
	postgres "tec-doc/internal/tec-doc/store/postgres"
	"tec-doc/pkg/model"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AddFromExcel mocks base method.
func (m *MockService) AddFromExcel(bodyData io.Reader, ctx *gin.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFromExcel", bodyData, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddFromExcel indicates an expected call of AddFromExcel.
func (mr *MockServiceMockRecorder) AddFromExcel(bodyData, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFromExcel", reflect.TypeOf((*MockService)(nil).AddFromExcel), bodyData, ctx)
}

// ExcelTemplateForClient mocks base method.
func (m *MockService) ExcelTemplateForClient() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExcelTemplateForClient")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExcelTemplateForClient indicates an expected call of ExcelTemplateForClient.
func (mr *MockServiceMockRecorder) ExcelTemplateForClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExcelTemplateForClient", reflect.TypeOf((*MockService)(nil).ExcelTemplateForClient))
}

// GetArticles mocks base method.
func (m *MockService) GetArticles(ctx context.Context, dataSupplierID int, article string) ([]model.Article, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetArticles", ctx, dataSupplierID, article)
	ret0, _ := ret[0].([]model.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetArticles indicates an expected call of GetArticles.
func (mr *MockServiceMockRecorder) GetArticles(ctx, dataSupplierID, article interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetArticles", reflect.TypeOf((*MockService)(nil).GetArticles), ctx, dataSupplierID, article)
}

// GetBrand mocks base method.
func (m *MockService) GetBrand(ctx context.Context, brandName string) (*model.Brand, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBrand", ctx, brandName)
	ret0, _ := ret[0].(*model.Brand)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBrand indicates an expected call of GetBrand.
func (mr *MockServiceMockRecorder) GetBrand(ctx, brandName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBrand", reflect.TypeOf((*MockService)(nil).GetBrand), ctx, brandName)
}

// GetProductsHistory mocks base method.
func (m *MockService) GetProductsHistory(ctx context.Context, tx postgres.Transaction, uploadID int64, limit, offset int) ([]model.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductsHistory", ctx, tx, uploadID, limit, offset)
	ret0, _ := ret[0].([]model.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductsHistory indicates an expected call of GetProductsHistory.
func (mr *MockServiceMockRecorder) GetProductsHistory(ctx, tx, uploadID, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductsHistory", reflect.TypeOf((*MockService)(nil).GetProductsHistory), ctx, tx, uploadID, limit, offset)
}

// GetSupplierTaskHistory mocks base method.
func (m *MockService) GetSupplierTaskHistory(ctx context.Context, tx postgres.Transaction, supplierID int64, limit, offset int) ([]model.Task, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSupplierTaskHistory", ctx, tx, supplierID, limit, offset)
	ret0, _ := ret[0].([]model.Task)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSupplierTaskHistory indicates an expected call of GetSupplierTaskHistory.
func (mr *MockServiceMockRecorder) GetSupplierTaskHistory(ctx, tx, supplierID, limit, offset interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSupplierTaskHistory", reflect.TypeOf((*MockService)(nil).GetSupplierTaskHistory), ctx, tx, supplierID, limit, offset)
}

// MockServer is a mock of Server interface.
type MockServer struct {
	ctrl     *gomock.Controller
	recorder *MockServerMockRecorder
}

// MockServerMockRecorder is the mock recorder for MockServer.
type MockServerMockRecorder struct {
	mock *MockServer
}

// NewMockServer creates a new mock instance.
func NewMockServer(ctrl *gomock.Controller) *MockServer {
	mock := &MockServer{ctrl: ctrl}
	mock.recorder = &MockServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockServer) EXPECT() *MockServerMockRecorder {
	return m.recorder
}

// Start mocks base method.
func (m *MockServer) Start() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start")
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockServerMockRecorder) Start() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockServer)(nil).Start))
}

// Stop mocks base method.
func (m *MockServer) Stop() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop")
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockServerMockRecorder) Stop() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockServer)(nil).Stop))
}

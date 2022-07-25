package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthorize(t *testing.T) {
	testCases := map[string]struct {
		header http.Header
		want   gin.H
	}{
		"success test": {
			header: http.Header{"X-User-Id": {"242"}, "X-Supplier-Id": {"1432"}},
			want:   gin.H{"X-User-Id": int64(242), "X-Supplier-Id": int64(1432)},
		},
		"success test with spaces": {
			header: http.Header{"X-User-Id": {"  242 "}, "X-Supplier-Id": {"  1432   "}},
			want:   gin.H{"X-User-Id": int64(242), "X-Supplier-Id": int64(1432)},
		},
		"invalid test with spaces": {
			header: http.Header{"X-User-Id": {"2 4 2"}, "X-Supplier-Id": {"1 43 2"}},
			want:   make(gin.H),
		},
		"invalid test": {
			header: http.Header{"X-User-Id": {"qwe"}, "X-Supplier-Id": {"qweqw"}},
			want:   make(gin.H),
		},
		"invalid User test": {
			header: http.Header{"X-User-Id": {"-4"}, "X-Supplier-Id": {"1432"}},
			want:   gin.H{"X-Supplier-Id": int64(1432)},
		},
		"invalid Supplier test": {
			header: http.Header{"X-User-Id": {"242"}, "X-Supplier-Id": {"blah foo bar"}},
			want:   gin.H{"X-User-Id": int64(242)},
		},
	}
	zerolog.SetGlobalLevel(zerolog.Disabled)
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c := &gin.Context{Request: &http.Request{Header: tc.header}}
			Authorize(c)

			wantUser, wantUserExist := tc.want["X-User-Id"]
			gotUser, gotUserExist := c.Get("X-User-Id")
			assert.Equal(t, wantUser, gotUser)
			assert.Equal(t, wantUserExist, gotUserExist)

			wantSupplier, wantSupplierExist := tc.want["X-Supplier-Id"]
			gotSupplier, gotSupplierExist := c.Get("X-Supplier-Id")
			assert.Equal(t, wantSupplier, gotSupplier)
			assert.Equal(t, wantSupplierExist, gotSupplierExist)
		})
	}
}

func TestCredentialsFromContext(t *testing.T) {
	type returnValues struct {
		supplierID int64
		userID     int64
		err        error
	}
	testCases := map[string]struct {
		input gin.H
		want  returnValues
	}{
		"success test": {
			input: gin.H{"X-User-Id": int64(22), "X-Supplier-Id": int64(22)},
			want:  returnValues{int64(22), int64(22), nil},
		},
		"invalid type test int16 User-id": {
			input: gin.H{"X-User-Id": int16(22), "X-Supplier-Id": int64(22)},
			want:  returnValues{0, 0, errors.New("user_id isn't type int64")},
		},
		"invalid type test int32 Supplier-id": {
			input: gin.H{"X-User-Id": int64(22), "X-Supplier-Id": int32(22)},
			want:  returnValues{0, 0, errors.New("supplier_id isn't type int64")},
		},
		"invalid type test (just int) User-id": {
			input: gin.H{"X-User-Id": 22, "X-Supplier-Id": int64(22)},
			want:  returnValues{0, 0, errors.New("user_id isn't type int64")},
		},
		"invalid type test (just int) Supplier-id": {
			input: gin.H{"X-User-Id": int64(22), "X-Supplier-Id": 22},
			want:  returnValues{0, 0, errors.New("supplier_id isn't type int64")},
		},
		"empty test user-id": {
			input: gin.H{"X-Supplier-Id": int64(22)},
			want:  returnValues{0, 0, errors.New("can't get user_id from context")},
		},
		"empty test supplier-id": {
			input: gin.H{"X-User-Id": int64(22)},
			want:  returnValues{0, 0, errors.New("can't get supplier_id from context")},
		},
	}

	gin.SetMode(gin.TestMode)
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(new(httptest.ResponseRecorder))
			c.Keys = tc.input

			supplierID, userID, err := CredentialsFromContext(c)
			assert.Equal(t, tc.want.supplierID, supplierID)
			assert.Equal(t, tc.want.userID, userID)
			assert.Equal(t, tc.want.err, err)
		})
	}
}

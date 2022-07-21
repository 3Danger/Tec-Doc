package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"strconv"
)

func Authorize(next *gin.Context) {
	userID := next.Request.Header.Get("X-User-Id")
	if userID != "" {
		userIDN, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			log.Error().Err(err).Str("Authorize", err.Error()).Send()
		} else {
			next.Set("X-User-Id", userIDN)
		}
	}

	supplierID := next.Request.Header.Get("X-Supplier-Id")
	if supplierID != "" {
		supplierIDN, err := strconv.ParseInt(supplierID, 10, 64)
		if err != nil {
			log.Error().Err(err).Str("Authorize", err.Error()).Send()
		} else {
			next.Set("X-Supplier-Id", supplierIDN)
		}
	}
}

func CredentialsFromContext(ctx *gin.Context) (supplierID int64, userID int64, err error) {
	var (
		valueUserID     interface{}
		valueSupplierID interface{}
	)
	if valueUserID = ctx.Value("X-User-Id"); valueUserID == nil {
		return 0, 0, fmt.Errorf("can't get user_id from context")
	}
	if valueSupplierID = ctx.Value("X-Supplier-Id"); valueSupplierID == nil {
		return 0, 0, fmt.Errorf("can't get supplier_id from context")
	}
	return valueUserID.(int64), valueSupplierID.(int64), nil
}

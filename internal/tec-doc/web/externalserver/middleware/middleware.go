package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

func Authorize(next *gin.Context) {
	userID := next.Request.Header.Get("X-User-Id")
	if userID != "" {
		userIDN, err := strconv.ParseInt(strings.TrimSpace(userID), 10, 64)
		if err != nil {
			log.Error().Err(err).Str("Authorize", err.Error()).Send()
		} else if userIDN >= 0 {
			next.Set("X-User-Id", userIDN)
		}
	}

	supplierID := next.Request.Header.Get("X-Supplier-Id")
	if supplierID != "" {
		supplierIDN, err := strconv.ParseInt(strings.TrimSpace(supplierID), 10, 64)
		if err != nil {
			log.Error().Err(err).Str("Authorize", err.Error()).Send()
		} else if supplierIDN >= 0 {
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
	valueUserIDN, ok := valueUserID.(int64)
	if !ok {
		return 0, 0, fmt.Errorf("user_id isn't type int64")
	}

	if valueSupplierID = ctx.Value("X-Supplier-Id"); valueSupplierID == nil {
		return 0, 0, fmt.Errorf("can't get supplier_id from context")
	}
	valueSupplierIDN, ok := valueSupplierID.(int64)
	if !ok {
		return 0, 0, fmt.Errorf("supplier_id isn't type int64")
	}
	return valueUserIDN, valueSupplierIDN, nil
}

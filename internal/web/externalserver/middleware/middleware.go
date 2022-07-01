package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"
)

func Authorize(next *gin.Context) {
	userID := next.Request.Header.Get("X-User-Id")
	if userID == "" {
		next.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id"})
		return
	}

	supplierID := next.Request.Header.Get("X-Supplier-Id")
	if userID == "" {
		next.JSON(http.StatusUnauthorized, gin.H{"error": "invalid supplier_id"})
		return
	}

	//TODO узнать как правильно добавить контексты
	userIDN, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("Authorize", err.Error()).Send()
	}
	supplierIDN, err := strconv.ParseInt(supplierID, 10, 64)
	if err != nil {
		log.Error().Err(err).Str("Authorize", err.Error()).Send()
	}
	next.Set("X-User-Id", userIDN)
	next.Set("X-Supplier-Id", supplierIDN)
	//next.ServeHTTP(w, req)
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

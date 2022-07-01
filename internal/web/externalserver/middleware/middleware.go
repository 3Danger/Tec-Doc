package middleware

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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
	ctx := next.Request.
		WithContext(context.WithValue(context.TODO(), "X-User-Id", userID)).
		WithContext(context.WithValue(context.TODO(), "X-Supplier-Id", supplierID))
	next.Request.WithContext(ctx.Context())
	next.Next()
	//next.ServeHTTP(w, req)
}

func CredentialsFromContext(ctx *gin.Context) (supplierID int64, userID int64, err error) {
	userID = ctx.Value("X-User-Id").(int64)
	if userID == 0 {
		return 0, 0, fmt.Errorf("can't get user_id from context")
	}

	supplierID = ctx.Value("X-Supplier-Id").(int64)
	if userID == 0 {
		return 0, 0, fmt.Errorf("can't get supplier_id from context")
	}

	return userID, supplierID, nil
}

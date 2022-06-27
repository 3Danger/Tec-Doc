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
	next.Request.
		WithContext(context.WithValue(context.TODO(), "X-User-Id", userID)).
		WithContext(context.WithValue(context.TODO(), "X-Supplier-Id", supplierID))

	next.Next()
	//next.ServeHTTP(w, req)
}

func СredentialsFromContext(ctx *gin.Context) (string, string, error) {
	userID := ctx.Value("X-User-Id").(string)
	if userID == "" {
		return "", "", fmt.Errorf("can't get user_id from context")
	}

	supplierID := ctx.Value("X-Supplier-Id").(string)
	if userID == "" {
		return "", "", fmt.Errorf("can't get supplier_id from context")
	}

	return userID, supplierID, nil
}

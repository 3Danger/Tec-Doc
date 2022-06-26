package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Printf("handling auth at %s\n", req.URL.Path)

		userID := req.Header.Get("X-User-Id")
		if userID == "" {
			http.Error(w, "invalid user_id", http.StatusUnauthorized)
				return
		}

		supplierID := req.Header.Get("X-Supplier-Id")
		if userID == "" {
			http.Error(w, "invalid supplier_id", http.StatusUnauthorized)
				return
		}

		ctx := context.WithValue(context.TODO(), "X-User-Id", userID)
		ctx = context.WithValue(ctx, "X-Supplier-Id", supplierID)
		req = req.WithContext(ctx)

		next.ServeHTTP(w, req)
	})
}

func СredentialsFromContext(ctx context.Context) (string, string, error) {
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

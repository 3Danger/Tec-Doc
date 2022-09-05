package externalserver

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"
	"strings"
	"tec-doc/pkg/errinfo"
	m "tec-doc/pkg/metrics"
	"time"
)

func (e *externalHttpServer) Authorize(ctx *gin.Context) {

	var (
		feature string
		scope   = e.service.Scope()
	)

	switch ctx.Request.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		feature = e.service.Scope().UpdateContentFeatureKey
	case http.MethodGet:
		feature = scope.ContentFeatureKey
	}

	userID := ctx.Request.Header.Get("X-User-Id")
	if userID == "" {
		e.logger.Error().Err(errinfo.InvalidUserID).Send()
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	userIDN, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		e.logger.Error().Err(errinfo.InvalidUserID).Send()
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	supplierID := ctx.Request.Header.Get("X-Supplier-Id")
	if supplierID == "" {
		e.logger.Error().Err(errinfo.InvalidSupplierID).Send()
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	supplierIDN, err := strconv.ParseInt(supplierID, 10, 64)
	if err != nil {
		e.logger.Error().Err(errinfo.InvalidSupplierID).Send()
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if supplierIDN == 0 || userIDN == 0 {
		e.logger.Error().Err(errinfo.InvalidUserOrSupplierID).Send()
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	key, err := uuid.FromString(supplierID)
	if err != nil {
		e.logger.Error().Err(errinfo.SupplierIsNotUUID).Send()
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	//TODO oldID
	_, err = e.service.Suppliers().GetOldSupplierID(ctx, key)
	if err != nil {
		e.logger.Error().Err(errinfo.FailOldSupplierID).Send()
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var decision bool
	if decision, err = e.service.Abac().CheckAccess(ctx, e.service.Scope().Scope, feature, &userIDN, key); err != nil {
		e.logger.Error().Err(errinfo.CheckAcessError).Send()
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if !decision {
		e.logger.Error().Err(errinfo.CheckAcessError).Send()
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Set("X-User-Id", userIDN)
	ctx.Set("X-Supplier-Id", supplierIDN)
}

func CredentialsFromContext(ctx *gin.Context) (supplierID, userID int64, err error) {
	valueUserID := ctx.GetInt64("X-User-Id")
	if valueUserID == 0 {
		return 0, 0, errinfo.InvalidUserID
	}

	valueSupplierID := ctx.GetInt64("X-Supplier-Id")
	if valueSupplierID == 0 {
		return 0, 0, errinfo.InvalidSupplierID
	}
	return valueUserID, valueSupplierID, nil
}

func (e *externalHttpServer) MiddleWareMetric(c *gin.Context) {
	t := time.Now()
	c.Next()
	status := strconv.Itoa(c.Writer.Status())
	e.metrics.Collector.WithLabelValues(
		m.InternalServerComponent,
		c.Request.Method,
		c.Request.URL.Path,
		status,
	).Inc()

	defer func() {
		e.metrics.LeadTime.WithLabelValues(
			m.InternalServerComponent,
			c.Request.Method,
			c.Request.URL.Path,
			strconv.FormatInt(time.Since(t).Milliseconds(), 10),
		).Observe(float64(time.Since(t).Milliseconds()))
	}()

	defer func() {
		e.metrics.LeadTimeQua.WithLabelValues(
			m.InternalServerComponent,
			c.Request.Method,
			c.Request.URL.Path,
			strconv.FormatInt(time.Since(t).Milliseconds(), 10),
		).Observe(float64(time.Since(t).Milliseconds()))
	}()

	e.metrics.Rating.WithLabelValues(
		m.InternalServerComponent,
		c.Request.Method,
		c.Request.URL.Path,
		status,
	).Inc()
}

func (e *externalHttpServer) corsMiddleware(ctx *gin.Context) {
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", strings.Join([]string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions,
	}, ", "))
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-Supplier-Id, X-User-Id")

	if ctx.Request.Method == http.MethodOptions {
		ctx.AbortWithStatus(204)
		return
	}
}

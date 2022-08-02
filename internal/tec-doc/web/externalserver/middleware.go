package externalserver

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"tec-doc/pkg/errinfo"
	m "tec-doc/pkg/metrics"
	"time"
)

func (e *externalHttpServer) Authorize(ctx *gin.Context) {
	userID := ctx.Request.Header.Get("X-User-Id")
	if userID == "" {
		e.logger.Error().Err(errinfo.InvalidUserID).Send()
		ctx.AbortWithStatus(401)
		return
	}

	userIDN, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		e.logger.Error().Err(errinfo.InvalidUserID).Send()
		ctx.AbortWithStatus(401)
		return
	}

	supplierID := ctx.Request.Header.Get("X-Supplier-Id")
	if supplierID == "" {
		e.logger.Error().Err(errinfo.InvalidSupplierID).Send()
		ctx.AbortWithStatus(401)
		return
	}

	supplierIDN, err := strconv.ParseInt(supplierID, 10, 64)
	if err != nil {
		e.logger.Error().Err(errinfo.InvalidSupplierID).Send()
		ctx.AbortWithStatus(401)
		return
	}
	if supplierIDN == 0 || userIDN == 0 {
		e.logger.Error().Err(errinfo.InvalidUserOrSupplierID).Send()
		ctx.AbortWithStatus(401)
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

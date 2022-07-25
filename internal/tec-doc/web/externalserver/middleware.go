package externalserver

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"strconv"
	m "tec-doc/pkg/metrics"
	"time"
)

// todo: если пришли пустые хедеры авторизации - возвращаем ошибку
// todo: именуем ошибки русскими аналогами понятными пользователю
// todo: сделать общую модель респонса
// todo: использовать логер структуры externalHttpServer
func (e *externalHttpServer) Authorize(next *gin.Context) {
	userID := next.Request.Header.Get("X-User-Id")
	if userID != "" {
		userIDN, err := strconv.ParseInt(userID, 10, 64)
		if err != nil {
			e.logger.Error().Err(err).Send()
		} else if userIDN >= 0 {
			next.Set("X-User-Id", userIDN)
		}
	}

	supplierID := next.Request.Header.Get("X-Supplier-Id")
	if supplierID != "" {
		supplierIDN, err := strconv.ParseInt(supplierID, 10, 64)
		if err != nil {
			log.Error().Err(err).Str("Authorize", err.Error()).Send()
		} else if supplierIDN >= 0 {
			next.Set("X-Supplier-Id", supplierIDN)
		}
	}
}

func CredentialsFromContext(ctx *gin.Context) (supplierID, userID int64, err error) {
	valueUserID := ctx.GetInt64("X-User-Id")
	if valueUserID == 0 {
		return 0, 0, fmt.Errorf("can't get user_id from context")
	}

	valueSupplierID := ctx.GetInt64("X-Supplier-Id")
	if valueSupplierID == 0 {
		return 0, 0, fmt.Errorf("can't get supplier_id from context")
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

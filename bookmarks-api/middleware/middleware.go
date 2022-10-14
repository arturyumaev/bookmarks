package middleware

import (
	"strconv"

	"github.com/arturyumaev/bookmarks/bookmarks-api/metrics"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

func PrometheusMiddleware(ctx *gin.Context) {
	path := ctx.Request.URL.Path
	timer := prometheus.NewTimer(metrics.HttpDuration.WithLabelValues(path))

	ctx.Next()

	statusCode := ctx.Writer.Status()
	metrics.ResponseStatus.WithLabelValues(strconv.Itoa(statusCode)).Inc()
	metrics.TotalRequests.WithLabelValues(path).Inc()

	timer.ObserveDuration()
}

package server

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Config struct {
	MetricPath        string
	MetricMaxRequests int
}

func NewServer(cfg *Config) *gin.Engine {

	perfEventCollector, _ := NewPerfEventCollector()

	rg := prometheus.NewRegistry()
	rg.Register(perfEventCollector)
	promHandler := promhttp.HandlerFor(
		rg,
		promhttp.HandlerOpts{
			MaxRequestsInFlight: cfg.MetricMaxRequests,
		},
	)

	r := gin.Default()

	r.GET(cfg.MetricPath, gin.WrapH(promHandler))

	v1 := r.Group("/api/v1")
	v1.Group("/collector/perfevent").
		POST("", addPerfEventCollector).
		POST("/del", delPerfEventCollector)

	return r
}

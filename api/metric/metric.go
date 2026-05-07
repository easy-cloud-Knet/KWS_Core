package metric

import (
	"net/http"

	"github.com/easy-cloud-Knet/KWS_Core/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (H *Handler) DefaultMetric() http.Handler {
	metricRegister := metrics.Initialize()

	return promhttp.HandlerFor(metricRegister, promhttp.HandlerOpts{})
}

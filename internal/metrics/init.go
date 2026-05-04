package metrics

import (
	"github.com/easy-cloud-Knet/KWS_Core/internal/metrics/ping"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

type Metrics interface {
	Enroll(prometheus.Registerer) error
}

var metricsList []Metrics = []Metrics{
	&ping.Collector{},
}

func Initialize() *prometheus.Registry {
	reg := prometheus.NewRegistry()

	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)
	for _, m := range metricsList {
		if err := m.Enroll(reg); err != nil {
			panic(err)
		}
	}
	return reg
}

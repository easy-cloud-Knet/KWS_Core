package ping

import (
	"time"

	probing "github.com/prometheus-community/pro-bing"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	rtt      *prometheus.GaugeVec
	dropRate *prometheus.GaugeVec
}

func (c *Collector) Enroll(reg prometheus.Registerer) error {
	c.rtt = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ping_rtt_seconds",
		Help: "Average round trip time in seconds",
	}, []string{"destination"})

	c.dropRate = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ping_drop_rate",
		Help: "Packet loss rate (0.0-100.0)",
	}, []string{"destination"})

	if err := reg.Register(c.rtt); err != nil {
		return err
	}
	if err := reg.Register(c.dropRate); err != nil {
		return err
	}

	go c.collect()
	return nil
}

func (c *Collector) collect() {
	for {
		c.ping("8.8.8.8")
		time.Sleep(10 * time.Second)
	}
}

func (c *Collector) ping(target string) {
	pinger, err := probing.NewPinger(target)
	if err != nil {
		return
	}
	pinger.Count = 3
	if err := pinger.Run(); err != nil {
		return
	}
	stats := pinger.Statistics()
	c.rtt.WithLabelValues(target).Set(stats.AvgRtt.Seconds())
	c.dropRate.WithLabelValues(target).Set(stats.PacketLoss)
}

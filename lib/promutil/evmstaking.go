package promutil

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	EVMStakingQueueDepth = promauto.NewGauge(prometheus.GaugeOpts{ //nolint:promlinter // skip
		Name: "evmstaking_queue_depth",
	})
)

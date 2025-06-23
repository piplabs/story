package promutil

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	EVMEngineExecutionPayloadMsgSize = promauto.NewGauge(prometheus.GaugeOpts{ //nolint:promlinter // skip
		Name: "evmengine_execution_payload_msg_size",
	})
)

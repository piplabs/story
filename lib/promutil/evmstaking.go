package promutil

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	EVMStakingWithdrawalQueueDepth = promauto.NewGauge(prometheus.GaugeOpts{ //nolint:promlinter // skip
		Name: "evmstaking_withdrawal_queue_depth",
	})
	EVMStakingRewardQueueDepth = promauto.NewGauge(prometheus.GaugeOpts{ //nolint:promlinter // skip
		Name: "evmstaking_reward_queue_depth",
	})
)

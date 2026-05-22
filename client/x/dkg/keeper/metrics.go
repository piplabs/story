package keeper

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Label values for roundsTotal "result" label.
const (
	labelRoundInitiated = "initiated"
	labelRoundCompleted = "completed"
	labelRoundSkipped   = "skipped"
)

// Label values for kernelCallDuration/Total "operation" label.
const (
	labelOpGenerateAndSealKey    = "generate_and_seal_key"
	labelOpGenerateDeals         = "generate_deals"
	labelOpProcessDeals          = "process_deals"
	labelOpProcessResponses      = "process_responses"
	labelOpProcessJustifications = "process_justifications"
	labelOpFinalizeDKG           = "finalize_dkg"
	labelOpPartialDecryptTDH2    = "partial_decrypt_tdh2"
)

// Label values for decryptRequestTotal "result" label.
const (
	labelDecryptStaleDropped = "stale_dropped"
	labelDecryptKernelFailed = "kernel_failed"
	labelDecryptRequeued     = "requeued"
	labelDecryptSubmitted    = "submitted"
)

// Label values for decryptBatchTotal "result" label.
const (
	labelBatchSuccess = "success"
	labelBatchError   = "error"
)

// Label values for pendingDataTotal "type" label.
const (
	labelPendingDeals          = "deals"
	labelPendingResponses      = "responses"
	labelPendingJustifications = "justifications"
)

// Label values for pendingDataTotal "op" label.
const (
	labelPendingCached   = "cached"
	labelPendingReplayed = "replayed"
	labelPendingDropped  = "dropped"
)

// Label values for kernelClientLookupDuration "result" label.
const (
	labelLookupHit       = "hit"
	labelLookupReconnect = "reconnect"
	labelLookupError     = "error"
)

var (
	roundsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "story",
		Subsystem: "dkg",
		Name:      "rounds_total",
		Help:      "Total number of DKG rounds by result (initiated, completed, skipped)",
	}, []string{"result"})

	kernelCallDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "story",
		Subsystem: "dkg",
		Name:      "kernel_call_duration_seconds",
		Help:      "Duration of kernel gRPC calls in seconds by operation",
		Buckets:   prometheus.DefBuckets,
	}, []string{"operation"})

	kernelCallTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "story",
		Subsystem: "dkg",
		Name:      "kernel_call_total",
		Help:      "Total number of kernel gRPC calls by operation and result",
	}, []string{"operation", "result"})

	decryptRequestTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "story",
		Subsystem: "dkg",
		Name:      "decrypt_request_total",
		Help:      "Total number of decrypt requests by result (stale_dropped, kernel_failed, requeued, submitted)",
	}, []string{"result"})

	decryptBatchTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "story",
		Subsystem: "dkg",
		Name:      "decrypt_batch_total",
		Help:      "Total number of decrypt batch submissions by result",
	}, []string{"result"})

	decryptBatchSize = promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "story",
		Subsystem: "dkg",
		Name:      "decrypt_batch_size",
		Help:      "Number of partial decryptions per successful batch submission",
		Buckets:   []float64{1, 5, 10, 15, 20, 25},
	})

	pendingDataTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "story",
		Subsystem: "dkg",
		Name:      "pending_data_total",
		Help:      "Total DKG pending data items by type (deals, responses, justifications) and op (cached, replayed)",
	}, []string{"type", "op"})

	sessionRecoveryTotal = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "story",
		Subsystem: "dkg",
		Name:      "session_recovery_total",
		Help:      "Total number of failed DKG sessions dispatched for recovery",
	})

	kernelClientLookupDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "story",
		Subsystem: "dkg",
		Name:      "kernel_client_lookup_duration_seconds",
		Help:      "Duration of getClientWithReconnect by result (hit, reconnect, error)",
		Buckets:   prometheus.DefBuckets,
	}, []string{"result"})
)

// observeKernelCall records duration and result for a single kernel gRPC operation.
// Call it immediately after the retry block resolves, passing the start time and
// the error returned by retry (nil on success).
func observeKernelCall(operation string, start time.Time, err error) {
	result := "success"
	if err != nil {
		result = "error"
	}

	kernelCallDuration.WithLabelValues(operation).Observe(time.Since(start).Seconds())
	kernelCallTotal.WithLabelValues(operation, result).Inc()
}

// incDecryptRequest adds n to the decrypt_request_total counter for result.
// result values: "stale_dropped", "kernel_failed", "requeued", "submitted".
func incDecryptRequest(result string, n int) {
	decryptRequestTotal.WithLabelValues(result).Add(float64(n))
}

// incDecryptBatch records a batch submission attempt.
// On success it also observes the batch size histogram.
func incDecryptBatch(result string, size int) {
	decryptBatchTotal.WithLabelValues(result).Inc()
	if result == "success" {
		decryptBatchSize.Observe(float64(size))
	}
}

// incPendingData adds n to the pending_data_total counter for (typ, op).
// typ: "deals", "responses", "justifications". op: "cached", "replayed".
func incPendingData(typ, op string, n int) {
	pendingDataTotal.WithLabelValues(typ, op).Add(float64(n))
}

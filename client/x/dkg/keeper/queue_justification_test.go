package keeper

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/piplabs/story/client/x/dkg/types"
)

// makeTestJustification creates a non-nil Justification pointer for testing.
func makeTestJustification(index uint32) *types.Justification {
	return &types.Justification{
		Index: index,
		VssJustification: &types.VSSJustification{
			SessionId: []byte("session"),
			Index:     index,
			PlainDeal: &types.PlainDeal{
				SessionId: []byte("pd-session"),
				SecShare: &types.SecShare{
					I: index,
					V: &types.Scalar{Data: []byte("scalar-data")},
				},
				Threshold: 3,
				Commitments: []*types.Point{
					{Data: []byte("commitment1")},
				},
			},
			Signature: []byte("sig"),
		},
	}
}

// drainJustifications clears the global justifications queue so tests are isolated.
func drainJustifications() {
	justificationsMu.Lock()
	defer justificationsMu.Unlock()
	justifications = nil
}

// TestEnqueueJustifications_ValidItems verifies that valid justifications are
// appended to the queue and can be dequeued in FIFO order.
func TestEnqueueJustifications_ValidItems(t *testing.T) {
	drainJustifications()
	defer drainJustifications()

	k, _ := setupDKGKeeper(t)

	j1 := makeTestJustification(1)
	j2 := makeTestJustification(2)
	j3 := makeTestJustification(3)

	k.EnqueueJustifications([]*types.Justification{j1, j2, j3})

	got := k.DequeueJustifications(3)
	require.Len(t, got, 3)
	require.Equal(t, uint32(1), got[0].Index)
	require.Equal(t, uint32(2), got[1].Index)
	require.Equal(t, uint32(3), got[2].Index)
}

// TestEnqueueJustifications_FiltersNilItems verifies that nil pointers in the
// input slice are silently skipped and not enqueued.
func TestEnqueueJustifications_FiltersNilItems(t *testing.T) {
	drainJustifications()
	defer drainJustifications()

	k, _ := setupDKGKeeper(t)

	j1 := makeTestJustification(10)
	j3 := makeTestJustification(30)

	// Pass a slice with nils interspersed
	k.EnqueueJustifications([]*types.Justification{j1, nil, j3, nil})

	got := k.DequeueJustifications(10)
	require.Len(t, got, 2, "nil entries must be filtered out")
	require.Equal(t, uint32(10), got[0].Index)
	require.Equal(t, uint32(30), got[1].Index)
}

// TestEnqueueJustifications_AllNil verifies that enqueueing a slice of all nils
// results in an empty queue.
func TestEnqueueJustifications_AllNil(t *testing.T) {
	drainJustifications()
	defer drainJustifications()

	k, _ := setupDKGKeeper(t)

	k.EnqueueJustifications([]*types.Justification{nil, nil, nil})

	got := k.DequeueJustifications(10)
	require.Nil(t, got, "all-nil slice should produce empty queue")
}

// TestDequeueJustifications_ReturnsUpToCount verifies that DequeueJustifications
// returns exactly the requested number when there are more items than count.
func TestDequeueJustifications_ReturnsUpToCount(t *testing.T) {
	drainJustifications()
	defer drainJustifications()

	k, _ := setupDKGKeeper(t)

	for i := uint32(1); i <= 5; i++ {
		k.EnqueueJustifications([]*types.Justification{makeTestJustification(i)})
	}

	// Dequeue only 2 of the 5
	first := k.DequeueJustifications(2)
	require.Len(t, first, 2)
	require.Equal(t, uint32(1), first[0].Index)
	require.Equal(t, uint32(2), first[1].Index)

	// Remaining 3 should still be available
	rest := k.DequeueJustifications(10)
	require.Len(t, rest, 3)
	require.Equal(t, uint32(3), rest[0].Index)
}

// TestDequeueJustifications_EmptyQueueReturnsNil verifies that dequeuing from
// an empty queue returns nil.
func TestDequeueJustifications_EmptyQueueReturnsNil(t *testing.T) {
	drainJustifications()
	defer drainJustifications()

	k, _ := setupDKGKeeper(t)

	got := k.DequeueJustifications(5)
	require.Nil(t, got, "empty queue should return nil")
}

// TestDequeueJustifications_CountExceedsQueueLength verifies that if count is
// larger than the queue size, all available items are returned.
func TestDequeueJustifications_CountExceedsQueueLength(t *testing.T) {
	drainJustifications()
	defer drainJustifications()

	k, _ := setupDKGKeeper(t)

	k.EnqueueJustifications([]*types.Justification{
		makeTestJustification(7),
		makeTestJustification(8),
	})

	got := k.DequeueJustifications(100)
	require.Len(t, got, 2, "should return all items when count > queue length")
}

// TestEnqueueDequeueJustifications_FIFO verifies that justifications are dequeued
// in the same order they were enqueued (FIFO).
func TestEnqueueDequeueJustifications_FIFO(t *testing.T) {
	drainJustifications()
	defer drainJustifications()

	k, _ := setupDKGKeeper(t)

	// Enqueue in two batches
	k.EnqueueJustifications([]*types.Justification{
		makeTestJustification(100),
		makeTestJustification(200),
	})
	k.EnqueueJustifications([]*types.Justification{
		makeTestJustification(300),
	})

	got := k.DequeueJustifications(10)
	require.Len(t, got, 3)
	require.Equal(t, uint32(100), got[0].Index)
	require.Equal(t, uint32(200), got[1].Index)
	require.Equal(t, uint32(300), got[2].Index)
}

// TestEnqueueJustifications_Concurrent verifies that concurrent enqueue and
// dequeue operations do not panic or corrupt state (goroutine-safety check).
func TestEnqueueJustifications_Concurrent(t *testing.T) {
	drainJustifications()
	defer drainJustifications()

	k, _ := setupDKGKeeper(t)

	const goroutines = 10
	const itemsPerGoroutine = 5

	var wg sync.WaitGroup

	// Concurrent enqueues
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(base int) {
			defer wg.Done()
			batch := make([]*types.Justification, itemsPerGoroutine)
			for j := 0; j < itemsPerGoroutine; j++ {
				batch[j] = makeTestJustification(uint32(base*itemsPerGoroutine + j))
			}
			k.EnqueueJustifications(batch)
		}(i)
	}

	wg.Wait()

	// Concurrent dequeues: drain everything
	var mu sync.Mutex
	var totalDequeued int

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			got := k.DequeueJustifications(itemsPerGoroutine)
			mu.Lock()
			totalDequeued += len(got)
			mu.Unlock()
		}()
	}

	wg.Wait()

	// Ensure we did not double-dequeue items (total <= goroutines*itemsPerGoroutine)
	require.LessOrEqual(t, totalDequeued, goroutines*itemsPerGoroutine,
		"concurrent dequeues must not exceed total enqueued items")
}

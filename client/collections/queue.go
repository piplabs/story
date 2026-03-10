package collections

import (
	"context"
	"errors"

	"cosmossdk.io/collections"
	"cosmossdk.io/collections/codec"

	ierrors "github.com/piplabs/story/lib/errors"
)

var (
	ErrPeek             = errors.New("queue peek failed")
	ErrEmptyQueue       = errors.New("queue is empty")
	ErrOutOfBoundsQueue = errors.New("queue index is out of bounds")
)

const (
	QueueElementsNameSuffix   = "_elements"
	QueueFrontNameSuffix      = "_front"
	QueueRearNameSuffix       = "_rear"
	QueueElementsPrefixSuffix = 0x0
	QueueFrontPrefixSuffix    = 0x1
	QueueRearPrefixSuffix     = 0x2
)

func NewQueue[T any](sb *collections.SchemaBuilder, prefix collections.Prefix, name string, vc codec.ValueCodec[T]) Queue[T] {
	return Queue[T]{
		front:    collections.NewItem(sb, append(prefix, QueueFrontPrefixSuffix), name+QueueFrontNameSuffix, collections.Uint64Value),
		rear:     collections.NewItem(sb, append(prefix, QueueRearPrefixSuffix), name+QueueRearNameSuffix, collections.Uint64Value),
		elements: collections.NewMap(sb, append(prefix, QueueElementsPrefixSuffix), name+QueueElementsNameSuffix, collections.Uint64Key, vc),
	}
}

type Queue[T any] struct {
	front    collections.Item[uint64]
	rear     collections.Item[uint64]
	elements collections.Map[uint64, T]
}

func (q Queue[T]) Initialize(ctx context.Context) error {
	// initialize the front and rear
	err := q.front.Set(ctx, 0)
	if err != nil {
		return ierrors.Wrap(err, "initialize set front")
	}

	err = q.rear.Set(ctx, 0)
	if err != nil {
		return ierrors.Wrap(err, "initialize set rear")
	}

	return nil
}

func (q Queue[T]) Enqueue(ctx context.Context, elem T) error {
	rear, err := q.rear.Get(ctx)
	if err != nil && !errors.Is(err, collections.ErrNotFound) {
		return ierrors.Wrap(err, "enqueue rear get")
	}

	err = q.elements.Set(ctx, rear, elem)
	if err != nil {
		return ierrors.Wrap(err, "enqueue elements set")
	}

	err = q.rear.Set(ctx, rear+1)
	if err != nil {
		return ierrors.Wrap(err, "enqueue rear set")
	}

	return nil
}

func (q Queue[T]) Dequeue(ctx context.Context) (elem T, err error) {
	if q.IsEmpty(ctx) {
		return elem, ierrors.Wrap(ErrEmptyQueue, "dequeue")
	}

	front, _ := q.front.Get(ctx)

	elem, err = q.elements.Get(ctx, front)
	if err != nil {
		return elem, ierrors.Wrap(err, "dequeue elements get")
	}

	err = q.elements.Remove(ctx, front)
	if err != nil {
		return elem, ierrors.Wrap(err, "dequeue elements remove")
	}

	err = q.front.Set(ctx, front+1)
	if err != nil {
		return elem, ierrors.Wrap(err, "dequeue front set")
	}

	return elem, nil
}

func (q Queue[T]) Get(ctx context.Context, index uint64) (elem T, err error) {
	// adjust the index to start at the front of the elements
	front, _ := q.front.Get(ctx)

	rear, _ := q.rear.Get(ctx)
	if front+index >= rear {
		return elem, ierrors.Wrap(ErrOutOfBoundsQueue, "get")
	}

	elem, err = q.elements.Get(ctx, front+index)
	if err != nil {
		return elem, ierrors.Wrap(err, "get elements")
	}

	return elem, nil
}

func (q Queue[T]) Len(ctx context.Context) uint64 {
	front, _ := q.front.Get(ctx)
	rear, _ := q.rear.Get(ctx)

	return rear - front
}

func (q Queue[T]) IsEmpty(ctx context.Context) bool {
	return q.Len(ctx) == 0
}

func (q Queue[T]) Peek(ctx context.Context) (elem T, err error) {
	if q.IsEmpty(ctx) {
		return elem, ierrors.Wrap(ErrEmptyQueue, "peek")
	}

	front, err := q.front.Get(ctx)
	if err != nil {
		return elem, ierrors.Wrap(err, "peek front get")
	}

	elem, err = q.elements.Get(ctx, front)
	if err != nil {
		return elem, ierrors.Wrap(err, "peek elements get")
	}

	return elem, nil
}

func (q Queue[T]) Front(ctx context.Context) (uint64, error) {
	item, err := q.front.Get(ctx)
	if err != nil {
		return 0, ierrors.Wrap(err, "front get")
	}

	return item, nil
}

func (q Queue[T]) Rear(ctx context.Context) (uint64, error) {
	item, err := q.rear.Get(ctx)
	if err != nil {
		return 0, ierrors.Wrap(err, "rear get")
	}

	return item, nil
}

func (q Queue[T]) Iterate(ctx context.Context) (collections.Iterator[uint64, T], error) {
	front, _ := q.front.Get(ctx)
	rear, _ := q.rear.Get(ctx)

	iter, err := q.elements.Iterate(ctx, new(collections.Range[uint64]).StartInclusive(front).EndExclusive(rear))
	if err != nil {
		return collections.Iterator[uint64, T]{}, ierrors.Wrap(err, "iterate")
	}

	return iter, nil
}

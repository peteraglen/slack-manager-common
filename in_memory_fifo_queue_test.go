package common

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInMemoryFifoQueue(t *testing.T) {
	t.Run("full queue should produce timeout error", func(t *testing.T) {
		ctx := context.Background()
		queue := NewInMemoryFifoQueue(2, time.Millisecond)
		err := queue.Send(ctx, "C000000001", "dedupID_1", "body_1")
		assert.NoError(t, err)
		err = queue.Send(ctx, "C000000002", "dedupID_2", "body_2")
		assert.NoError(t, err)
		err = queue.Send(ctx, "C000000003", "dedupID_3", "body_3")
		assert.ErrorContains(t, err, "timeout")
	})

	t.Run("cancelled context should return context error", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		queue := NewInMemoryFifoQueue(1, time.Second)
		err := queue.Send(ctx, "C000000001", "dedupID_1", "body_1")
		assert.NoError(t, err)
		cancel()
		err = queue.Send(ctx, "C000000002", "dedupID_2", "body_2")
		assert.Error(t, err)
		assert.ErrorIs(t, err, context.Canceled)
	})

	t.Run("receive function should return all items in order", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		queue := NewInMemoryFifoQueue(3, time.Millisecond)
		err := queue.Send(ctx, "C000000001", "dedupID_1", "body_1")
		assert.NoError(t, err)
		err = queue.Send(ctx, "C000000002", "dedupID_2", "body_2")
		assert.NoError(t, err)
		err = queue.Send(ctx, "C000000003", "dedupID_3", "body_3")
		assert.NoError(t, err)

		receivedItems := make(chan *FifoQueueItem, 3)

		go func() {
			err := queue.Receive(ctx, receivedItems)
			assert.ErrorIs(t, err, context.Canceled)
		}()

		result := []*FifoQueueItem{}

		for item := range receivedItems {
			result = append(result, item)
			if len(result) == 3 {
				cancel()
			}
		}

		assert.Len(t, result, 3)
		assert.Equal(t, "body_1", result[0].Body)
		assert.Equal(t, "body_2", result[1].Body)
		assert.Equal(t, "body_3", result[2].Body)

		assert.NoError(t, result[0].Ack(context.Background()))
		assert.Nil(t, result[0].ExtendVisibility)
	})

	t.Run("receive function should react to context cancelled when waiting to write to sinkCh", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		queue := NewInMemoryFifoQueue(2, time.Second)
		err := queue.Send(ctx, "C000000001", "dedupID_1", "body_1")
		assert.NoError(t, err)
		err = queue.Send(ctx, "C000000002", "dedupID_2", "body_2")
		assert.NoError(t, err)

		receivedItems := make(chan *FifoQueueItem)

		go func() {
			err := queue.Receive(ctx, receivedItems)
			assert.ErrorIs(t, err, context.Canceled)
		}()

		for range receivedItems {
			cancel()
			break
		}
	})
}

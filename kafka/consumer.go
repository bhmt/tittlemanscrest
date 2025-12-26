package kafka

import (
	"context"
	"errors"
	"io"
	"log"
	"time"

	kgo "github.com/segmentio/kafka-go"
)

var DefaultReaderConfig = kgo.ReaderConfig{
	MaxWait:          1 * time.Second,
	ReadBatchTimeout: 1 * time.Second,
	CommitInterval:   1 * time.Second,
}

type Consumer[T any] struct {
	reader     *kgo.Reader
	serializer Serializer[T]
	backoffMin time.Duration
	backoffMax time.Duration
}

func WithBackoffMin[T any](n time.Duration) func(*Consumer[T]) error {
	return func(c *Consumer[T]) error {
		c.backoffMin = n
		return nil
	}
}

func WithBackoffMax[T any](n time.Duration) func(*Consumer[T]) error {
	return func(c *Consumer[T]) error {
		c.backoffMax = n
		return nil
	}
}

type HandlerFunc[T any] func(ctx context.Context, msg T) error

func NewConsumer[T any](brokers []string, topic, groupID string, serializer Serializer[T], config *kgo.ReaderConfig) *Consumer[T] {
	readerConfig := DefaultReaderConfig

	if config != nil {
		readerConfig = *config
	}

	readerConfig.Brokers = brokers
	readerConfig.Topic = topic
	readerConfig.GroupID = groupID

	return &Consumer[T]{
		reader:     kgo.NewReader(readerConfig),
		serializer: serializer,
		backoffMin: 100 * time.Millisecond,
		backoffMax: 1 * time.Second,
	}
}

func (c *Consumer[T]) Consume(ctx context.Context, handler HandlerFunc[T]) error {
	defer c.reader.Close()

	currentBackoff := c.backoffMin

	for {
		m, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			if errors.Is(err, io.EOF) {
				return nil
			}

			select {
			case <-ctx.Done():
				return nil
			case <-time.After(currentBackoff):
				currentBackoff *= 2
				if currentBackoff > c.backoffMax {
					currentBackoff = c.backoffMax
				}
			}
			continue
		}

		currentBackoff = c.backoffMin

		payload, err := c.serializer.Deserialize(m.Value)
		if err != nil {
			log.Printf("serialization error: %v\n", err)
			c.reader.CommitMessages(ctx, m)
			continue
		}

		if err := handler(ctx, payload); err != nil {
			log.Printf("handler failed: %v\n", err)
			// retry or DLQ.
			continue
		}

		if err := c.reader.CommitMessages(ctx, m); err != nil {
			log.Printf("failed to commit message: %v\n", err)
		}
	}
}

func (c *Consumer[T]) Close() {
	c.reader.Close()
}

package kafka

import (
	"context"
	"time"

	kgo "github.com/segmentio/kafka-go"
)

var DefaultWriter = kgo.Writer{
	WriteTimeout: 1 * time.Second,
	BatchTimeout: 1 * time.Second,
}

type Producer[T any] struct {
	writer     *kgo.Writer
	serializer Serializer[T]
}

func NewProducer[T any](brokers []string, topic string, serializer Serializer[T], writer *kgo.Writer) (*Producer[T], error) {
	defaultWriter := &DefaultWriter
	if writer != nil {
		defaultWriter = writer
	}

	defaultWriter.Addr = kgo.TCP(brokers...)
	defaultWriter.Topic = topic
	defaultWriter.Balancer = &kgo.LeastBytes{}

	p := &Producer[T]{
		writer:     defaultWriter,
		serializer: serializer,
	}

	return p, nil
}

func (p *Producer[T]) Publish(ctx context.Context, key string, msg T) error {
	data, err := p.serializer.Serialize(msg)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kgo.Message{
		Key:   []byte(key),
		Value: data,
	})
}

func (p *Producer[T]) Close() error {
	return p.writer.Close()
}

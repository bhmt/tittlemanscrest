//go:build integration

package kafka_test

import (
	"context"
	"testing"
	"time"

	"github.com/bhmt/tittlemanscrest/kafka"
	"github.com/stretchr/testify/assert"
)

var brokers = []string{"localhost:49092"}
var topic = "items"

type item struct {
	Id        int    `json:"id"`
	Timestamp int64  `json:"ts"`
	Value     string `json:"value"`
}

var testSerializer = kafka.JsonSerializer[item]{}

func TestIntegrationPublisherConsumer(t *testing.T) {
	// allow kafka and kafka setup to run
	time.Sleep(5 * time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	prod, err := kafka.NewProducer(brokers, topic, testSerializer, nil)
	assert.NoError(t, err, err)
	defer prod.Close()

	cons := kafka.NewConsumer(brokers, topic, "item-test-group", testSerializer, nil)
	defer cons.Close()

	testItem := item{
		Id:        1,
		Timestamp: time.Now().Unix(),
		Value:     "antigravity",
	}

	recieve := make(chan item)

	handler := func(ctx context.Context, msg item) error {
		recieve <- msg
		return nil
	}

	go func() {
		err := cons.Consume(ctx, handler)
		assert.NoError(t, err, err)
	}()

	go func() {
		err := prod.Publish(ctx, "item-test-key", testItem)
		assert.NoError(t, err, err)
	}()

	select {
	case msg := <-recieve:
		assert.Equal(t, testItem.Id, msg.Id)
		assert.Equal(t, testItem.Timestamp, msg.Timestamp)
		assert.Equal(t, testItem.Value, msg.Value)
		return
	case <-time.After(15 * time.Second):
		t.Fatal("test timeout exceded")
	}
}

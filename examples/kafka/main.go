package main

import (
	"context"
	"fmt"
	"log"
	"math/rand/v2"
	"time"

	"github.com/bhmt/tittlemanscrest/cmd"
	"github.com/bhmt/tittlemanscrest/kafka"
)

var brokers = []string{"localhost:49092"}
var topic = "items"

type Item struct {
	Id        int    `json:"id"`
	Timestamp int64  `json:"ts"`
	Value     string `json:"value"`
}

func produce(ctx context.Context) {
	prod, err := kafka.NewProducer(brokers, topic, kafka.JsonSerializer[Item]{}, nil)
	if err != nil {
		log.Printf("new producer error: %v", err)
		return
	}

	defer prod.Close()

	i := 0
	for {
		tick := time.Duration(rand.IntN(2)+2) * time.Second

		select {
		case <-ctx.Done():
			return
		case <-time.After(tick):
			key := fmt.Sprintf("item%d", i)
			item := Item{
				Id:        i,
				Timestamp: time.Now().Unix(),
				Value:     fmt.Sprintf("%d%s", int(tick), key),
			}

			if err := prod.Publish(ctx, key, item); err != nil {
				log.Printf("producer error: %v", err)
				return
			}
			i++
		}
	}
}

func consume(ctx context.Context) {
	cons := kafka.NewConsumer(brokers, topic, "item-service-group", kafka.JsonSerializer[Item]{}, nil)
	defer cons.Close()

	handler := func(ctx context.Context, msg Item) error {
		log.Printf("item %d %d %s", msg.Id, msg.Timestamp, msg.Value)
		return nil
	}

	if err := cons.Consume(ctx, handler); err != nil {
		log.Printf("consumer error: %v", err)
		return
	}
}

func main() {
	cmd.Run(func(ctx context.Context) {
		go produce(ctx)
		go consume(ctx)
		<-ctx.Done()
	})
}

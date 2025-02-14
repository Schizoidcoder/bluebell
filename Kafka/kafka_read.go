package Kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func InitKafkaReader(ctx context.Context, topic string) {
	Kafka_like_reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		Topic:          topic,
		CommitInterval: time.Second, //每个一段时间上报一次offset
		GroupID:        "test",      //每一个组消费一份topic
		StartOffset:    kafka.FirstOffset,
	})
	for {
		msg, err := Kafka_like_reader.ReadMessage(ctx)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			fmt.Println(string(msg.Value))
		}
	}
}

package Kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func InitKafkaReader(ctx context.Context, topics []string) {
	Kafka_reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"localhost:9092"},
		GroupTopics:    topics,
		CommitInterval: time.Second, //每个一段时间上报一次offset
		GroupID:        "test",      //每一个组消费一份topic
		StartOffset:    kafka.FirstOffset,
	})
	for {
		msg, err := Kafka_reader.ReadMessage(ctx)
		if err != nil {
			fmt.Println(err)
			break
		} else {
			fmt.Println(msg.Topic, string(msg.Key), string(msg.Value))
		}
	}
}

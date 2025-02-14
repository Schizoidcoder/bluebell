package Kafka

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"
)

func GetKafkaWriter(ctx context.Context, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP("localhost:9092"),
		Topic:                  topic,
		Balancer:               &kafka.Hash{}, //负载均衡
		WriteTimeout:           1 * time.Second,
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: true,
	}
}

package Kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/segmentio/kafka-go"
)

var (
	Kafka_like_writer *kafka.Writer
	Kafka_like_reader *kafka.Reader
)

func Init() {
	var wg sync.WaitGroup
	wg.Add(1)
	var ctx context.Context = context.Background()
	Kafka_like_writer = GetKafkaWriter(ctx, "like_event")
	defer Kafka_like_writer.Close()
	go func() {
		wg.Done()
		InitKafkaReader(ctx, "like_event")
	}()
	wg.Wait()
}

func Close() {
	if Kafka_like_reader != nil {
		err := Kafka_like_reader.Close()
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
	}
}

package Kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/segmentio/kafka-go"
)

var (
	Kafka_writer *kafka.Writer
	Kafka_reader *kafka.Reader
)

func Init() {
	var wg sync.WaitGroup
	wg.Add(1)
	var ctx context.Context = context.Background()
	Kafka_writer = GetKafkaWriter(ctx, "like_event")
	go func() {
		wg.Done()
		InitKafkaReader(ctx, []string{"like_event", "comment_event"})
	}()
	wg.Wait()
}

func Close() {
	if Kafka_writer != nil {
		err := Kafka_writer.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
	if Kafka_reader != nil {
		err := Kafka_reader.Close()
		if err != nil {
			fmt.Println(err)
		}
	}
}

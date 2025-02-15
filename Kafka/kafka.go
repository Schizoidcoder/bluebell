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
	Ctx          context.Context = context.Background()
)

func Init() {
	var wg sync.WaitGroup
	wg.Add(1)
	Kafka_writer = GetKafkaWriter(Ctx, "like_event")
	go func() {
		wg.Done()
		InitKafkaReader(Ctx, []string{"like_event", "comment_event"})
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

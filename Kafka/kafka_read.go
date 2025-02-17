package Kafka

import (
	"bluebell/models"
	mywebsocket "bluebell/websocket"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
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
			var event models.Event
			err = json.Unmarshal(msg.Value, &event)
			if err != nil {
				fmt.Println(err)
				continue
			}
			itoa := strconv.FormatInt(event.AuthorID, 10)
			AimClient, flag := mywebsocket.CheckIfConnected(itoa)
			if !flag {
				log.Printf("当前用户已下线")
				//Todo:用户下线后的逻辑
				continue
			}
			AimClient.Send <- msg.Value
			//fmt.Println(msg.Topic, string(msg.Key), event)
		}
	}
}

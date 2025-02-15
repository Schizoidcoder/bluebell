package Kafka

import (
	"bluebell/models"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

func KafkaSendMessage(userId int64, authorId int64, postId string, username string, topic string, content string) error {
	event := &models.Event{
		UserID: userId, AuthorID: authorId, PostID: postId, UserName: username, Action: content,
	}
	eventBytes, err := json.Marshal(event)
	if err != nil {
		zap.L().Error("KafkaSendMessage", zap.Error(err))
		return err
	}
	err = Kafka_writer.WriteMessages(Ctx, kafka.Message{
		Topic: topic,
		Key:   []byte(fmt.Sprintf("%d", userId)),
		Value: eventBytes,
	})
	if err != nil {
		zap.L().Error("KafkaSendMessage", zap.Error(err))
		return err
	}
	return nil
}

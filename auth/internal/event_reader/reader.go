package event_reader

import (
	global "async_course/main"
	"async_course/main/internal/service"
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/go-playground/validator"
	"github.com/segmentio/kafka-go"
)

type EventReader struct {
	s *service.Service
}

func NewEventReader(s *service.Service) *EventReader {
	return &EventReader{s}
}

func (er *EventReader) StartReaders(brokers []string, groupID string) {

	// topic A
	topicAReader := newTopicReader(brokers, groupID, global.KafkaTopicIDA)
	go handle(context.Background(), topicAReader, er.handleMessageJSON)

	// topic B
	topicBReader := newTopicReader(brokers, groupID, global.KafkaTopicIDB)
	go handle(context.Background(), topicBReader, er.handleMessageJSON)

}

func newTopicReader(brokers []string, groupID string, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   brokers,
		GroupID:   groupID,
		Topic:     topic,
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
}

func closeReader(r *kafka.Reader) {
	if err := r.Close(); err != nil {
		slog.Error("failed to close reader", "error", err)
		os.Exit(1)
	}
}

type messageHandler func(m kafka.Message) error

func handle(ctx context.Context, r *kafka.Reader, fn messageHandler) {
	defer closeReader(r)
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			slog.Error("error while reading message", "error", err)
			break
		}
		slog.Info("received message from broker", "topic", r.Config().Topic, "key", string(m.Key), "value", string(m.Value))
		if err := fn(m); err != nil {
			slog.Error("error while handling message", "error", err)
		}
	}
}

func validatePayload[T any](m kafka.Message) (T, error) {
	var payload T
	if err := json.Unmarshal(m.Value, &payload); err != nil {
		slog.Error("error while unmarshaling payload", "key", string(m.Key), "value", string(m.Value), "error", err)
		return payload, err
	}
	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		slog.Error("error while validating payload", "key", string(m.Key), "value", string(m.Value), "error", err)
		return payload, err
	}
	return payload, nil
}

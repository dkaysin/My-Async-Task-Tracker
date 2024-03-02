package event_reader

import (
	schema "async_course/schema_registry"
	"async_course/task/internal/service"

	"context"
	"log/slog"
	"os"

	"github.com/segmentio/kafka-go"
)

type EventReader struct {
	s              *service.Service
	SchemaRegistry *schema.SchemaRegistry
}

func NewEventReader(s *service.Service, sr *schema.SchemaRegistry) *EventReader {
	return &EventReader{s, sr}
}

func newTopicReader(brokers []string, groupID string, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 1,    // should be suitable for our case
		MaxBytes: 10e6, // 10MB
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
		slog.Info("received message from broker", "topic", r.Config().Topic)
		if err := fn(m); err != nil {
			slog.Error("error while handling message", "error", err)
		}
	}
}

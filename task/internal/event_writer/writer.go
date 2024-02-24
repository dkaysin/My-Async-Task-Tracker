package event_writer

import (
	"async_course/task"
	"context"
	"encoding/json"
	"log/slog"
	"os"

	"github.com/segmentio/kafka-go"
)

type EventWriter struct {
}

type TopicWriter struct {
}

func NewEventWriter(brokers []string) *EventWriter {
	return &EventWriter{
	}
}

func newTopicWriter(brokers []string, topic string) *TopicWriter {
	return &TopicWriter{&kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}}
}

func (er *EventWriter) Close() {
	if err := er.TopicWriterTask.w.Close(); err != nil {
		slog.Error("failed to close writer", "error", err)
		os.Exit(1)
	}
}

func (tr *TopicWriter) WriteBytes(ctx context.Context, key string, value []byte) error {
	err := tr.w.WriteMessages(ctx, kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	})
	if err != nil {
		slog.Error("failed to write message", "topic", tr.w.Topic, "key", key, "value", value, "error", err)
		return err
	}
	slog.Info("written message", "topic", tr.w.Topic, "key", key, "value", value)
	return nil
}

func (tr *TopicWriter) WriteString(ctx context.Context, key string, value string) error {
	return tr.WriteBytes(ctx, key, []byte(value))
}

func (tr *TopicWriter) WriteJSON(ctx context.Context, key string, value any) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		slog.Error("failed to marshall payload", "topic", tr.w.Topic, "key", key, "value", value, "error", err)
		return err
	}
	return tr.WriteBytes(ctx, key, valueBytes)
}

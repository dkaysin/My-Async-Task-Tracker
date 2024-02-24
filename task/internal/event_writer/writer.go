package event_writer

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

type TopicWriter struct {
	w *kafka.Writer
}

func NewEventWriter(brokers []string) *EventWriter {
	return &EventWriter{}
}

func newTopicWriter(brokers []string, topic string) *TopicWriter {
	return &TopicWriter{&kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}}
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

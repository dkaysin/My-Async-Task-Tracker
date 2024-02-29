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

func newTopicWriter(brokers []string, topic string) *TopicWriter {
	return &TopicWriter{&kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		BatchSize:    50, // should be suitable for our case
		BatchTimeout: 10, // should be suitable for our case
	}}
}

func (tr *TopicWriter) WriteBytes(key string, value []byte) error {
	ctx := context.Background()
	go func(ctx context.Context, key string, value []byte) {
		err := tr.w.WriteMessages(ctx, kafka.Message{
			Key:   []byte(key),
			Value: []byte(value),
		})
		if err != nil {
			slog.Error("failed to write message", "topic", tr.w.Topic, "key", key, "value", value, "error", err)
			return
		}
		slog.Info("written message", "topic", tr.w.Topic, "key", key, "value", value)

	}(ctx, key, value)
	return nil
}

func (tr *TopicWriter) WriteString(key string, value string) error {
	return tr.WriteBytes(key, []byte(value))
}

func (tr *TopicWriter) WriteJSON(key string, value any) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		slog.Error("failed to marshall payload", "topic", tr.w.Topic, "key", key, "value", value, "error", err)
		return err
	}
	return tr.WriteBytes(key, valueBytes)
}

package event_writer

import (
	schema "async_course/schema_registry"

	"context"
	"log/slog"
	"time"

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
		BatchSize:    50,                    // should be suitable for our case
		BatchTimeout: time.Millisecond * 10, // should be suitable for our case
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
			return
		}

	}(ctx, key, value)
	return nil
}

func (tr *TopicWriter) WriteMessage(m schema.Message) error {
	payloadBytes, err := schema.MarshalAndValidate(m.Event.PayloadSchema, m.Event.Payload)
	if err != nil {
		slog.Error("failed to marshall payload", "topic", tr.w.Topic, "key", m.Key, "event_name", m.Event.EventName, "event_version", m.Event.EventVersion, "error", err)
		return err
	}
	eventRaw := schema.EventRaw{
		EventName:     m.Event.EventName,
		EventID:       m.Event.EventID,
		EventVersion:  m.Event.EventVersion,
		EventProducer: m.Event.EventProducer,
		Payload:       payloadBytes,
	}
	eventBytes, err := schema.MarshalAndValidate(schema.EventRawSchema, eventRaw)
	if err != nil {
		slog.Error("failed to marshall event", "topic", tr.w.Topic, "key", m.Key, "event_name", m.Event.EventName, "event_version", m.Event.EventVersion, "error", err)
		return err
	}
	err = tr.WriteBytes(m.Key, eventBytes)
	if err != nil {
		slog.Error("error while writing message", "topic", tr.w.Topic, "key", m.Key, "event_name", m.Event.EventName, "event_version", m.Event.EventVersion, "error", err)
		return err
	}
	slog.Info("written message", "topic", tr.w.Topic, "key", m.Key, "event_name", m.Event.EventName, "event_version", m.Event.EventVersion)
	return nil
}

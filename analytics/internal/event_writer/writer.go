package event_writer

import (
	schema "async_course/schema_registry"
	general "async_course/schema_registry/schemas/general"

	"context"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/protocol"
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

func (tr *TopicWriter) WriteBytes(meta general.Meta, key string, value []byte) error {
	ctx := context.Background()
	go func(ctx context.Context, key string, value []byte) {
		err := tr.w.WriteMessages(ctx, kafka.Message{
			Key:   []byte(key),
			Value: value,
			Headers: []protocol.Header{
				{
					Key:   "event_name",
					Value: []byte(meta.EventName),
				},
				{
					Key:   "event_version",
					Value: []byte(meta.EventVersion),
				},
				{
					Key:   "event_producer",
					Value: []byte(meta.EventProducer),
				},
				{
					Key:   "event_id",
					Value: []byte(meta.EventID),
				},
			},
		})
		if err != nil {
			return
		}

	}(ctx, key, value)
	return nil
}

func (tr *TopicWriter) WriteMessage(m schema.Message) error {
	payloadBytes, err := schema.MarshalAndValidate(m.PayloadSchema, m.Payload)
	if err != nil {
		slog.Error("failed to marshall payload", "topic", tr.w.Topic, "key", m.Key, "event_name", m.Meta.EventName, "event_version", m.Meta.EventVersion, "error", err)
		return err
	}
	err = tr.WriteBytes(m.Meta, m.Key, payloadBytes)
	if err != nil {
		slog.Error("error while writing message", "topic", tr.w.Topic, "key", m.Key, "event_name", m.Meta.EventName, "event_version", m.Meta.EventVersion, "error", err)
		return err
	}
	slog.Info("written message", "topic", tr.w.Topic, "key", m.Key, "event_name", m.Meta.EventName, "event_version", m.Meta.EventVersion)
	return nil
}

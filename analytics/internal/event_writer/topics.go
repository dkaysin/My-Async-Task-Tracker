package event_writer

import (
	"async_course/analytics"
	schema "async_course/schema_registry"
	"log/slog"
	"os"
)

type EventWriter struct {
	TopicWriterPayment *TopicWriter
	SchemaRegistry     *schema.SchemaRegistry
}

func NewEventWriter(brokers []string, sr *schema.SchemaRegistry) *EventWriter {
	return &EventWriter{
		TopicWriterPayment: newTopicWriter(brokers, analytics.KafkaTopicPayment),
		SchemaRegistry:     sr,
	}
}

func (er *EventWriter) Close() {
	if err := er.TopicWriterPayment.w.Close(); err != nil {
		slog.Error("failed to close writer", "error", err)
		os.Exit(1)
	}
}

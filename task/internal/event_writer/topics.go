package event_writer

import (
	schema "async_course/schema_registry"
	"async_course/task"

	"log/slog"
	"os"
)

type EventWriter struct {
	TopicWriterTask *TopicWriter
	SchemaRegistry  *schema.SchemaRegistry
}

func NewEventWriter(brokers []string, sr *schema.SchemaRegistry) *EventWriter {
	return &EventWriter{
		TopicWriterTask: newTopicWriter(brokers, task.KafkaTopicTask),
		SchemaRegistry:  sr,
	}
}

func (er *EventWriter) Close() {
	if err := er.TopicWriterTask.w.Close(); err != nil {
		slog.Error("failed to close writer", "error", err)
		os.Exit(1)
	}
}

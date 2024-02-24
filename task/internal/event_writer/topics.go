package event_writer

import (
	"async_course/task"
	"log/slog"
	"os"
)

type EventWriter struct {
	TopicWriterTask *TopicWriter
}

func NewEventWriter(brokers []string) *EventWriter {
	return &EventWriter{
		TopicWriterTask: newTopicWriter(brokers, task.KafkaTopicTask),
	}
}
func (er *EventWriter) Close() {
	if err := er.TopicWriterTask.w.Close(); err != nil {
		slog.Error("failed to close writer", "error", err)
		os.Exit(1)
	}
}

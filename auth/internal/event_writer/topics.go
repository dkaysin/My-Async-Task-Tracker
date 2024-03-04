package event_writer

import (
	"async_course/auth"
	"log/slog"
	"os"
)

type EventWriter struct {
	TopicWriterAccount *TopicWriter
}

func NewEventWriter(brokers []string) *EventWriter {
	return &EventWriter{
		TopicWriterAccount: newTopicWriter(brokers, auth.KafkaTopicAccount),
	}
}

func (er *EventWriter) Close() {
	if err := er.TopicWriterAccount.w.Close(); err != nil {
		slog.Error("failed to close writer", "error", err)
		os.Exit(1)
	}
}

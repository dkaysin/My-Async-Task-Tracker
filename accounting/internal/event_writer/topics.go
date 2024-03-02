package event_writer

import (
	"async_course/accounting"
	schema "async_course/schema_registry"
)

type EventWriter struct {
	TopicWriterAccount *TopicWriter
	TopicWriterTask    *TopicWriter
	SchemaRegistry     *schema.SchemaRegistry
}

func NewEventWriter(brokers []string, sr *schema.SchemaRegistry) *EventWriter {
	return &EventWriter{
		TopicWriterAccount: newTopicWriter(brokers, accounting.KafkaTopicAccount),
		TopicWriterTask:    newTopicWriter(brokers, accounting.KafkaTopicTask),
		SchemaRegistry:     sr,
	}
}

func (er *EventWriter) Close() {
}

package event_writer

import "async_course/auth"

type EventWriter struct {
	TopicAWriter *TopicWriter
	TopicBWriter *TopicWriter
}

func NewEventWriter(brokers []string) *EventWriter {
	return &EventWriter{
		TopicAWriter: newTopicWriter(brokers, auth.KafkaTopicIDA),
		TopicBWriter: newTopicWriter(brokers, auth.KafkaTopicIDB),
	}
}

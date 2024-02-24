package event_writer

import "async_course/task"

type EventWriter struct {
	TopicWriterTask *TopicWriter
}

func NewEventWriter(brokers []string) *EventWriter {
	return &EventWriter{
		TopicWriterTask: newTopicWriter(brokers, task.KafkaTopicTask),
	}
}

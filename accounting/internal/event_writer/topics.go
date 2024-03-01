package event_writer

import (
	schema "async_course/schema_registry"
)

type EventWriter struct {
	SchemaRegistry *schema.SchemaRegistry
}

func NewEventWriter(brokers []string, sr *schema.SchemaRegistry) *EventWriter {
	return &EventWriter{
		SchemaRegistry: sr,
	}
}

func (er *EventWriter) Close() {
}

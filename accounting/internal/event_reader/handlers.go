package event_reader

import (
	schema "async_course/schema_registry"

	"log/slog"

	"github.com/segmentio/kafka-go"
)

func (er *EventReader) StartReaders(brokers []string, groupID string) {
}

func (er *EventReader) handleMessage(m kafka.Message) error {
	var eventRaw schema.EventRaw
	err := schema.UnmarshalAndValidate(schema.EventSchema, m.Value, &eventRaw)
	if err != nil {
		slog.Error("errorw while unmarshaling event", "err", err)
		return err
	}

	switch eventRaw.EventName {
	}
	return err
}

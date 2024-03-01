package event_reader

import (
	schema "async_course/schema_registry"
	"async_course/task"

	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

func (er *EventReader) StartReaders(brokers []string, groupID string) {
	topicReaderAccount := newTopicReader(brokers, groupID, task.KafkaTopicAccount)
	go handle(context.Background(), topicReaderAccount, er.handleMessage)
}

func (er *EventReader) handleMessage(m kafka.Message) error {
	var eventRaw schema.EventRaw
	err := schema.UnmarshalAndValidate(schema.EventSchema, m.Value, &eventRaw)
	if err != nil {
		slog.Error("errorw while unmarshaling event", "err", err)
		return err
	}

	switch eventRaw.EventName {
	case schema.EventNameAccountCreated:
		err = er.handleAccountCreated(eventRaw)
	case schema.EventNameAccountUpdated:
		err = er.handleAccountUpdated(eventRaw)
	}
	return err
}

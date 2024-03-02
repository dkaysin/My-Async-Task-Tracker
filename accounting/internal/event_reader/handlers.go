package event_reader

import (
	"async_course/accounting"
	schema "async_course/schema_registry"

	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

func (er *EventReader) StartReaders(brokers []string, groupID string) {
	topicReaderAccount := newTopicReader(brokers, groupID, accounting.KafkaTopicAccount)
	go handle(context.Background(), topicReaderAccount, er.handleMessage)

	topicReaderTask := newTopicReader(brokers, groupID, accounting.KafkaTopicTask)
	go handle(context.Background(), topicReaderTask, er.handleMessage)
}

func (er *EventReader) handleMessage(m kafka.Message) error {
	var eventRaw schema.EventRaw
	err := schema.UnmarshalAndValidate(schema.EventRawSchema, m.Value, &eventRaw)
	if err != nil {
		slog.Error("errorw while unmarshaling event", "err", err)
		return err
	}

	switch eventRaw.EventName {
	case schema.EventNameAccountCreated:
		err = er.handleAccountCreated(eventRaw)
	case schema.EventNameAccountUpdated:
		err = er.handleAccountUpdated(eventRaw)
	case schema.EventNameTaskAssigned:
		err = er.handleTaskAssigned(eventRaw)
	case schema.EventNameTaskCompleted:
		err = er.handleTaskCompleted(eventRaw)
	}
	return err
}

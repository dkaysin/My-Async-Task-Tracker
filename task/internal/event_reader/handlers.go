package event_reader

import (
	schema "async_course/schema_registry"
	"async_course/task"

	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

func (er *EventReader) StartReaders(brokers []string, groupID string) {
	topicReaderAccount := newTopicReader(brokers, groupID, schema.KafkaTopicAccount)
	go handle(context.Background(), topicReaderAccount, er.handleMessage)
}

func (er *EventReader) handleMessage(m kafka.Message) error {
	eventName := getHeader(m, "event_name")
	if eventName == "" {
		return task.ErrMessageHeaderNotFound
	}
	eventVersion := getHeader(m, "event_version")
	if eventVersion == "" {
		return task.ErrMessageHeaderNotFound
	}
	slog.Info("received kafka message", "event_name", eventName, "event_version", eventVersion)

	var err error
	switch eventName {
	case schema.EventNameAccountCreated:
		switch eventVersion {
		case "1":
			err = er.handleAccountCreatedV1(m.Value)
		default:
			err = task.ErrUnknownEventVersion
		}
	case schema.EventNameAccountUpdated:
		switch eventVersion {
		case "1":
			err = er.handleAccountUpdatedV1(m.Value)
		default:
			err = task.ErrUnknownEventVersion
		}
	}
	if err != nil {
		slog.Error("error while handling message", "error", err)
	}
	return err
}

func getHeader(m kafka.Message, key string) string {
	for _, header := range m.Headers {
		if header.Key == key {
			return string(header.Value)
		}
	}
	slog.Error("header not found in kafka message", "key", key)
	return ""
}

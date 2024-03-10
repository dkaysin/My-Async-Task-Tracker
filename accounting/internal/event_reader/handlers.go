package event_reader

import (
	"async_course/accounting"
	schema "async_course/schema_registry"

	"context"
	"log/slog"

	"github.com/segmentio/kafka-go"
)

func (er *EventReader) StartReaders(brokers []string, groupID string) {
	topicReaderTask := newTopicReader(brokers, groupID, schema.KafkaTopicTask)
	go handle(context.Background(), topicReaderTask, er.handleMessage)

	topicReaderAccount := newTopicReader(brokers, groupID, schema.KafkaTopicAccount)
	go handle(context.Background(), topicReaderAccount, er.handleMessage)
}

func (er *EventReader) handleMessage(m kafka.Message) error {
	eventName := getHeader(m, "event_name")
	if eventName == "" {
		return accounting.ErrMessageHeaderNotFound
	}
	eventVersion := getHeader(m, "event_version")
	if eventVersion == "" {
		return accounting.ErrMessageHeaderNotFound
	}
	slog.Info("received kafka message", "event_name", eventName, "event_version", eventVersion)

	var err error
	switch eventName {
	case schema.EventNameTaskAssigned:
		switch eventVersion {
		case "1":
			err = er.handleTaskAssigned(m.Value)
		default:
			err = accounting.ErrUnknownEventVersion
		}
	case schema.EventNameTaskCompleted:
		switch eventVersion {
		case "1":
			err = er.handleTaskCompleted(m.Value)
		default:
			err = accounting.ErrUnknownEventVersion
		}
	case schema.EventNameAccountCreated:
		switch eventVersion {
		case "1":
			err = er.handleAccountCreatedV1(m.Value)
		default:
			err = accounting.ErrUnknownEventVersion
		}
	case schema.EventNameAccountUpdated:
		switch eventVersion {
		case "1":
			err = er.handleAccountUpdatedV1(m.Value)
		default:
			err = accounting.ErrUnknownEventVersion
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

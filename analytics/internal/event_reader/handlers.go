package event_reader

import (
	"async_course/analytics"
	schema "async_course/schema_registry"
	"context"

	"log/slog"

	"github.com/segmentio/kafka-go"
)

func (er *EventReader) StartReaders(brokers []string, groupID string) {
	topicReaderTransaction := newTopicReader(brokers, groupID, schema.KafkaTopicTransaction)
	go handle(context.Background(), topicReaderTransaction, er.handleMessage)

	topicReaderAccount := newTopicReader(brokers, groupID, schema.KafkaTopicAccount)
	go handle(context.Background(), topicReaderAccount, er.handleMessage)

	topicReaderTask := newTopicReader(brokers, groupID, schema.KafkaTopicTask)
	go handle(context.Background(), topicReaderTask, er.handleMessage)
}

func (er *EventReader) handleMessage(m kafka.Message) error {
	eventName := getHeader(m, "event_name")
	if eventName == "" {
		return analytics.ErrMessageHeaderNotFound
	}
	eventVersion := getHeader(m, "event_version")
	if eventVersion == "" {
		return analytics.ErrMessageHeaderNotFound
	}
	slog.Info("received kafka message", "event_name", eventName, "event_version", eventVersion)

	var err error
	switch eventName {
	case schema.EventNameTransactionRevenue:
		switch eventVersion {
		case "1":
			err = er.handleTransactionRevenueV1(m.Value)
		default:
			err = analytics.ErrUnknownEventVersion
		}
	case schema.EventNameTransactionCost:
		switch eventVersion {
		case "1":
			err = er.handleTransactionCostV1(m.Value)
		default:
			err = analytics.ErrUnknownEventVersion
		}
	case schema.EventNameTaskAssigned:
		switch eventVersion {
		case "1":
			err = er.handleTaskUpdatedV1(m.Value)
		case "2":
			err = er.handleTaskUpdatedV2(m.Value)
		default:
			err = analytics.ErrUnknownEventVersion
		}
	case schema.EventNameTaskCompleted:
		switch eventVersion {
		case "1":
			err = er.handleTaskUpdatedV1(m.Value)
		case "2":
			err = er.handleTaskUpdatedV2(m.Value)
		default:
			err = analytics.ErrUnknownEventVersion
		}
	case schema.EventNameAccountCreated:
		switch eventVersion {
		case "1":
			err = er.handleAccountCreatedV1(m.Value)
		default:
			err = analytics.ErrUnknownEventVersion
		}
	case schema.EventNameAccountUpdated:
		switch eventVersion {
		case "1":
			err = er.handleAccountUpdatedV1(m.Value)
		default:
			err = analytics.ErrUnknownEventVersion
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

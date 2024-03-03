package event_reader

import (
	"async_course/analytics"
	schema "async_course/schema_registry"
	"context"

	"log/slog"

	"github.com/segmentio/kafka-go"
)

func (er *EventReader) StartReaders(brokers []string, groupID string) {
	topicReaderTransaction := newTopicReader(brokers, groupID, analytics.KafkaTopicTransaction)
	go handle(context.Background(), topicReaderTransaction, er.handleMessage)
}

func (er *EventReader) handleMessage(m kafka.Message) error {
	var eventRaw schema.EventRaw
	err := schema.UnmarshalAndValidate(schema.EventRawSchema, m.Value, &eventRaw)
	if err != nil {
		slog.Error("errorw while unmarshaling event", "err", err)
		return err
	}
	slog.Info("parsed raw event", "event_name", eventRaw.EventName, "event_version", eventRaw.EventVersion, "event_producer", eventRaw.EventProducer)

	switch eventRaw.EventName {
	case schema.EventNameTransactionRevenue:
		err = er.handleTransactionRevenue(eventRaw)
	case schema.EventNameTransactionCost:
		err = er.handleTransactionCost(eventRaw)
	}
	if err != nil {
		slog.Error("error while handling message", "error", err)
	}
	return err
}

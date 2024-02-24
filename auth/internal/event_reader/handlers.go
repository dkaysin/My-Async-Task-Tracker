package event_reader

import (
	"async_course/auth"
	"context"

	"github.com/segmentio/kafka-go"
)

func (er *EventReader) StartReaders(brokers []string, groupID string) {
	topicAReader := newTopicReader(brokers, groupID, auth.KafkaTopicIDA)
	go handle(context.Background(), topicAReader, er.handleMessageJSON)
}

func (er *EventReader) handleMessageJSON(m kafka.Message) error {
	switch string(m.Key) {
	// handle incoming events
	}
	return nil
}

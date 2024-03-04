package event_reader

import (
	"github.com/segmentio/kafka-go"
)

func (er *EventReader) StartReaders(brokers []string, groupID string) {
}

func (er *EventReader) handleMessageJSON(m kafka.Message) error {
	switch string(m.Key) {
	}
	return nil
}

package event_reader

import (
	"github.com/segmentio/kafka-go"
)

func (er *EventReader) handleMessageJSON(m kafka.Message) error {
	var err error
	switch string(m.Key) {
	}
	return err
}

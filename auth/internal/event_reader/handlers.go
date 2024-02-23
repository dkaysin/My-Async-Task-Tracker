package event_reader

import (
	"async_course/auth"

	"github.com/segmentio/kafka-go"
)

func (er *EventReader) handleMessageJSON(m kafka.Message) error {
	switch string(m.Key) {
	case auth.Event1:
		// handle Event1
	case auth.Event2:
		// handle Event2
	}
	return nil
}

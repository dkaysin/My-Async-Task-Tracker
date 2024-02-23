package event_reader

import (
	global "async_course/main"

	"github.com/segmentio/kafka-go"
)

func (er *EventReader) handleMessageJSON(m kafka.Message) error {
	switch string(m.Key) {
	case global.Event1:
		// handle Event1
	case global.Event2:
		// handle Event2
	}
	return nil
}

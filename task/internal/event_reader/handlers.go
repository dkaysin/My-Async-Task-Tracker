package event_reader

import (
	"async_course/task"

	"github.com/segmentio/kafka-go"
)

func (er *EventReader) handleMessageJSON(m kafka.Message) error {
	var err error
	switch string(m.Key) {
	case task.EventKeyAccountCreated:
		err = er.handleAccountCreated(m)
	case task.EventKeyAccountUpdated:
		err = er.handleAccountUpdated(m)
	}
	return err
}

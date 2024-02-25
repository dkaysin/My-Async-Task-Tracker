package event_reader

import (
	"async_course/task"
	"context"

	"github.com/segmentio/kafka-go"
)

func (er *EventReader) handleAccountCreated(m kafka.Message) error {
	payload, err := validatePayload[task.EventValueAccountCreated](m)
	if err != nil {
		return err
	}
	return er.s.UpsertAccountRole(context.Background(), payload.UserID, true, payload.Role)
}

func (er *EventReader) handleAccountUpdated(m kafka.Message) error {
	payload, err := validatePayload[task.EventValueAccountUpdated](m)
	if err != nil {
		return err
	}
	return er.s.UpsertAccountRole(context.Background(), payload.UserID, payload.Active, payload.Role)
}

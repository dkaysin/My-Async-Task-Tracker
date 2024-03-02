package event_reader

import (
	schema "async_course/schema_registry"
	"log/slog"

	"context"
)

func (er *EventReader) handleTaskAssigned(e schema.EventRaw) error {
	var payload schema.EventPayloadTaskAssigned
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.TaskAssignedSchema, e.Payload, &payload)
	if err != nil {
		return err
	}
	if payload.OldUserID == &payload.Task.UserID {
		slog.Info("task owner did not change, skipping processing assign")
		return nil
	}
	return er.s.ProcessAssignTask(context.Background(), payload.Task.UserID)
}

func (er *EventReader) handleTaskCompleted(e schema.EventRaw) error {
	var payload schema.EventPayloadTaskCompleted
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.TaskCompletedSchema, e.Payload, &payload)
	if err != nil {
		return err
	}
	return er.s.ProcessCompleteTask(context.Background(), payload.Task.UserID)
}

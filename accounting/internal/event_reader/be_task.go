package event_reader

import (
	schema "async_course/schema_registry"

	"context"
)

func (er *EventReader) handleTaskAssigned(e schema.EventRaw) error {
	var payload schema.TaskAssigned
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.TaskAssignedSchema, e.Payload, &payload)
	if err != nil {
		return err
	}
	return er.s.ProcessAssignTask(context.Background(), payload.Task.UserID, payload.Task.TaskID)
}

func (er *EventReader) handleTaskCompleted(e schema.EventRaw) error {
	var payload schema.TaskCompleted
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.TaskCompletedSchema, e.Payload, &payload)
	if err != nil {
		return err
	}
	return er.s.ProcessCompleteTask(context.Background(), payload.Task.UserID, payload.Task.TaskID)
}

package event_reader

import (
	schema "async_course/schema_registry"

	"context"
)

func (er *EventReader) handleTaskAssigned(payload []byte) error {
	var event schema.TaskAssigned
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.V1.TaskAssignedSchema, payload, &event)
	if err != nil {
		return err
	}
	err = er.s.UpsertTask(context.Background(), event.Task)
	if err != nil {
		return err
	}
	return er.s.ProcessAssignTask(context.Background(), event.Task.UserID, event.Task.TaskID)
}

func (er *EventReader) handleTaskCompleted(payload []byte) error {
	var event schema.TaskCompleted
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.V1.TaskCompletedSchema, payload, &event)
	if err != nil {
		return err
	}
	err = er.s.UpsertTask(context.Background(), event.Task)
	if err != nil {
		return err
	}
	return er.s.ProcessCompleteTask(context.Background(), event.Task.UserID, event.Task.TaskID)
}

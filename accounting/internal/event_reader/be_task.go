package event_reader

import (
	schema "async_course/schema_registry"
	v1 "async_course/schema_registry/schemas/v1"
	v2 "async_course/schema_registry/schemas/v2"

	"context"
)

func (er *EventReader) handleTaskAssignedV1(payload []byte) error {
	var event v1.TaskAssigned
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.V1.TaskAssignedSchema, payload, &event)
	if err != nil {
		return err
	}
	err = er.s.UpsertTaskV1(context.Background(), event.Task)
	if err != nil {
		return err
	}
	return er.s.ProcessAssignTask(context.Background(), event.Task.UserID, event.Task.TaskID)
}

func (er *EventReader) handleTaskCompletedV1(payload []byte) error {
	var event v1.TaskCompleted
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.V1.TaskCompletedSchema, payload, &event)
	if err != nil {
		return err
	}
	err = er.s.UpsertTaskV1(context.Background(), event.Task)
	if err != nil {
		return err
	}
	return er.s.ProcessCompleteTask(context.Background(), event.Task.UserID, event.Task.TaskID)
}

func (er *EventReader) handleTaskAssignedV2(payload []byte) error {
	var event v2.TaskAssigned
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.V2.TaskAssignedSchema, payload, &event)
	if err != nil {
		return err
	}
	err = er.s.UpsertTaskV2(context.Background(), event.Task)
	if err != nil {
		return err
	}
	return er.s.ProcessAssignTask(context.Background(), event.Task.UserID, event.Task.TaskID)
}

func (er *EventReader) handleTaskCompletedV2(payload []byte) error {
	var event v2.TaskCompleted
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.V1.TaskCompletedSchema, payload, &event)
	if err != nil {
		return err
	}
	err = er.s.UpsertTaskV2(context.Background(), event.Task)
	if err != nil {
		return err
	}
	return er.s.ProcessCompleteTask(context.Background(), event.Task.UserID, event.Task.TaskID)
}

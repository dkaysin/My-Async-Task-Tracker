package event_reader

import (
	schema "async_course/schema_registry"
	v1 "async_course/schema_registry/schemas/v1"
	v2 "async_course/schema_registry/schemas/v2"

	"context"
)

func (er *EventReader) handleTaskUpdatedV1(payload []byte) error {
	var event v1.TaskAssigned
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.V1.TaskAssignedSchema, payload, &event)
	if err != nil {
		return err
	}
	return er.s.UpsertTaskV1(context.Background(), event.Task)
}

func (er *EventReader) handleTaskUpdatedV2(payload []byte) error {
	var event v2.TaskAssigned
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.V2.TaskAssignedSchema, payload, &event)
	if err != nil {
		return err
	}
	return er.s.UpsertTaskV2(context.Background(), event.Task)
}

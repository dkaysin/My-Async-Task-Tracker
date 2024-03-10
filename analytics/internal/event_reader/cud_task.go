package event_reader

import (
	schema "async_course/schema_registry"
	"context"
)

func (er *EventReader) handleTaskUpdatedV1(payload []byte) error {
	var event schema.TaskAssigned
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.V1.TaskAssignedSchema, payload, &event)
	if err != nil {
		return err
	}
	return er.s.UpsertTask(context.Background(), event.Task)
}

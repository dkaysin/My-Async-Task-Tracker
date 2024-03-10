package event_reader

import (
	schema "async_course/schema_registry"
	v1 "async_course/schema_registry/schemas/v1"

	"context"
)

func (er *EventReader) handleAccountCreatedV1(payload []byte) error {
	var event v1.AccountCreated
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.V1.AccountCreatedSchema, payload, &event)
	if err != nil {
		return err
	}
	return er.s.UpsertAccountRole(context.Background(), event.UserID, true, event.Role)
}

func (er *EventReader) handleAccountUpdatedV1(payload []byte) error {
	var event v1.AccountUpdated
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.V1.AccountUpdatedSchema, payload, &event)
	if err != nil {
		return err
	}
	return er.s.UpsertAccountRole(context.Background(), event.UserID, event.Active, event.Role)
}

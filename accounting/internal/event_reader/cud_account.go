package event_reader

import (
	schema "async_course/schema_registry"

	"context"
)

func (er *EventReader) handleAccountCreated(e schema.EventRaw) error {
	var payload schema.AccountCreated
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.AccountCreatedSchema, e.Payload, &payload)
	if err != nil {
		return err
	}
	return er.s.UpsertAccountRole(context.Background(), payload.UserID, true, payload.Role)
}

func (er *EventReader) handleAccountUpdated(e schema.EventRaw) error {
	var payload schema.AccountUpdated
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.AccountUpdatedSchema, e.Payload, &payload)
	if err != nil {
		return err
	}
	return er.s.UpsertAccountRole(context.Background(), payload.UserID, payload.Active, payload.Role)
}

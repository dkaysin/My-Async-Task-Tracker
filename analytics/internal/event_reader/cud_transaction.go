package event_reader

import (
	schema "async_course/schema_registry"

	"context"
)

func (er *EventReader) handleTransactionRevenue(e schema.EventRaw) error {
	var payload schema.TransactionRevenue
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.TransactionRevenueSchema, e.Payload, &payload)
	if err != nil {
		return err
	}
	return er.s.ProcessTaskAssigned(context.Background(), payload.UserID, payload.Source, payload.Revenue, payload.CreatedAt)
}

func (er *EventReader) handleTransactionCost(e schema.EventRaw) error {
	var payload schema.TransactionCost
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.TransactionCostSchema, e.Payload, &payload)
	if err != nil {
		return err
	}
	return er.s.ProcessTaskCompleted(context.Background(), payload.UserID, payload.Source, payload.Cost, payload.CreatedAt)
}

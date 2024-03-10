package event_reader

import (
	schema "async_course/schema_registry"

	"context"
)

func (er *EventReader) handleTransactionRevenueV1(payload []byte) error {
	var event schema.TransactionRevenue
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.V1.TransactionRevenueSchema, payload, &event)
	if err != nil {
		return err
	}
	return er.s.ProcessTaskAssigned(context.Background(), event.UserID, event.Source, event.Revenue, event.CreatedAt)
}

func (er *EventReader) handleTransactionCostV1(payload []byte) error {
	var event schema.TransactionCost
	err := schema.UnmarshalAndValidate(er.SchemaRegistry.V1.TransactionCostSchema, payload, &event)
	if err != nil {
		return err
	}
	return er.s.ProcessTaskCompleted(context.Background(), event.UserID, event.Source, event.Cost, event.CreatedAt)
}

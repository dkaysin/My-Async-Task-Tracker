package schema_registry

import (
	general "async_course/schema_registry/schemas/general"
	v1 "async_course/schema_registry/schemas/v1"
	v2 "async_course/schema_registry/schemas/v2"
	"time"

	"github.com/google/uuid"
	"github.com/hamba/avro/v2"
)

type Message struct {
	Key           string
	Meta          general.Meta
	Payload       interface{} `avro:"payload"`
	PayloadSchema avro.Schema
}

// Task.Assigned

const EventNameTaskAssigned = "Task.Assigned"

func (s *SchemaV1) NewEventTaskAssigned(task v1.Task, oldUserID *string) Message {
	return Message{
		Key: task.TaskID,
		Meta: general.Meta{
			EventName:     EventNameTaskAssigned,
			EventID:       uuid.New().String(),
			EventVersion:  s.Version,
			EventProducer: s.Producer,
		},
		Payload: v1.TaskAssigned{
			Task:      task,
			OldUserID: oldUserID,
		},
		PayloadSchema: s.TaskAssignedSchema,
	}
}

func (s *SchemaV2) NewEventTaskAssigned(task v2.Task, oldUserID *string) Message {
	return Message{
		Key: task.TaskID,
		Meta: general.Meta{
			EventName:     EventNameTaskAssigned,
			EventID:       uuid.New().String(),
			EventVersion:  s.Version,
			EventProducer: s.Producer,
		},
		Payload: v2.TaskAssigned{
			Task:      task,
			OldUserID: oldUserID,
		},
		PayloadSchema: s.TaskAssignedSchema,
	}
}

// Task.Completed

const EventNameTaskCompleted = "Task.Completed"

func (s *SchemaV1) NewEventTaskCompleted(task v1.Task) Message {
	return Message{
		Key: task.TaskID,
		Meta: general.Meta{
			EventName:     EventNameTaskCompleted,
			EventID:       uuid.New().String(),
			EventVersion:  s.Version,
			EventProducer: s.Producer,
		},
		Payload: v1.TaskCompleted{
			Task: task,
		},
		PayloadSchema: s.TaskCompletedSchema,
	}
}

func (s *SchemaV2) NewEventTaskCompleted(task v2.Task) Message {
	return Message{
		Key: task.TaskID,
		Meta: general.Meta{
			EventName:     EventNameTaskCompleted,
			EventID:       uuid.New().String(),
			EventVersion:  s.Version,
			EventProducer: s.Producer,
		},
		Payload: v2.TaskCompleted{
			Task: task,
		},
		PayloadSchema: s.TaskCompletedSchema,
	}
}

// Account.Created

const EventNameAccountCreated = "Account.Created"

func (s *SchemaV1) NewEventAccountCreated(userID, role string) Message {
	return Message{
		Key: userID,
		Meta: general.Meta{
			EventName:     EventNameAccountCreated,
			EventID:       uuid.New().String(),
			EventVersion:  s.Version,
			EventProducer: s.Producer,
		},
		Payload: v1.AccountCreated{
			UserID: userID,
			Role:   role,
		},
		PayloadSchema: s.AccountCreatedSchema,
	}
}

// Account.Updated

const EventNameAccountUpdated = "Account.Updated"

func (s *SchemaV1) NewEventAccountUpdated(userID, role string, active bool) Message {
	return Message{
		Key: userID,
		Meta: general.Meta{
			EventName:     EventNameAccountCreated,
			EventID:       uuid.New().String(),
			EventVersion:  s.Version,
			EventProducer: s.Producer,
		},
		Payload: v1.AccountUpdated{
			UserID: userID,
			Role:   role,
			Active: active,
		},
		PayloadSchema: s.AccountCreatedSchema,
	}
}

// Payment.Made

const EventNamePaymentMade = "Payment.Made"

func (s *SchemaV1) NewEventPaymentMade(userID string, amount int, processedAt time.Time) Message {
	return Message{
		Key: userID,
		Meta: general.Meta{
			EventName:     EventNamePaymentMade,
			EventID:       uuid.New().String(),
			EventVersion:  s.Version,
			EventProducer: s.Producer,
		},
		Payload: v1.PaymentMade{
			UserID:      userID,
			Amount:      amount,
			ProcessedAt: processedAt,
		},
		PayloadSchema: s.PaymentMadeSchema,
	}
}

// Transaction.Profit

const EventNameTransactionRevenue = "Transaction.Revenue"

func (s *SchemaV1) NewEventTransactionRevenue(userID *string, taskID string, revenue int, createdAt time.Time) Message {
	var key string
	if userID != nil {
		key = *userID
	}
	return Message{
		Key: key,
		Meta: general.Meta{
			EventName:     EventNameTransactionRevenue,
			EventID:       uuid.New().String(),
			EventVersion:  s.Version,
			EventProducer: s.Producer,
		},
		Payload: v1.TransactionRevenue{
			UserID:    userID,
			Source:    taskID,
			Revenue:   revenue,
			CreatedAt: createdAt,
		},
		PayloadSchema: s.TransactionRevenueSchema,
	}
}

const EventNameTransactionCost = "Transaction.Cost"

func (s *SchemaV1) NewEventTransactionCost(userID *string, taskID string, cost int, createdAt time.Time) Message {
	var key string
	if userID != nil {
		key = *userID
	}
	return Message{
		Key: key,
		Meta: general.Meta{
			EventName:     EventNameTransactionCost,
			EventID:       uuid.New().String(),
			EventVersion:  s.Version,
			EventProducer: s.Producer,
		},
		Payload: v1.TransactionCost{
			UserID:    userID,
			Source:    taskID,
			Cost:      cost,
			CreatedAt: createdAt,
		},
		PayloadSchema: s.TransactionCostSchema,
	}
}

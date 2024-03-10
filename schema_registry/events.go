package schema_registry

import (
	"time"

	"github.com/google/uuid"
	"github.com/hamba/avro/v2"
)

type Message struct {
	Key           string
	Meta          Meta
	Payload       interface{} `avro:"payload"`
	PayloadSchema avro.Schema
}

// Task.Assigned

const EventNameTaskAssigned = "Task.Assigned"

func (s *SchemaV1) NewEventTaskAssigned(task Task, oldUserID *string) Message {
	return Message{
		Key: task.TaskID,
		Meta: Meta{
			EventName:     EventNameTaskAssigned,
			EventID:       uuid.New().String(),
			EventVersion:  s.Version,
			EventProducer: s.Producer,
		},
		Payload: TaskAssigned{
			Task:      task,
			OldUserID: oldUserID,
		},
		PayloadSchema: s.TaskAssignedSchema,
	}
}

// Task.Completed

const EventNameTaskCompleted = "Task.Completed"

func (s *SchemaV1) NewEventTaskCompleted(task Task) Message {
	return Message{
		Key: task.TaskID,
		Meta: Meta{
			EventName:     EventNameTaskCompleted,
			EventID:       uuid.New().String(),
			EventVersion:  s.Version,
			EventProducer: s.Producer,
		},
		Payload: TaskCompleted{
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
		Meta: Meta{
			EventName:     EventNameAccountCreated,
			EventID:       uuid.New().String(),
			EventVersion:  s.Version,
			EventProducer: s.Producer,
		},
		Payload: AccountCreated{
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
		Meta: Meta{
			EventName:     EventNameAccountCreated,
			EventID:       uuid.New().String(),
			EventVersion:  s.Version,
			EventProducer: s.Producer,
		},
		Payload: AccountUpdated{
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
		Meta: Meta{
			EventName:     EventNamePaymentMade,
			EventID:       uuid.New().String(),
			EventVersion:  s.Version,
			EventProducer: s.Producer,
		},
		Payload: PaymentMade{
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
		Meta: Meta{
			EventName:     EventNameTransactionRevenue,
			EventID:       uuid.New().String(),
			EventVersion:  s.Version,
			EventProducer: s.Producer,
		},
		Payload: TransactionRevenue{
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
		Meta: Meta{
			EventName:     EventNameTransactionCost,
			EventID:       uuid.New().String(),
			EventVersion:  s.Version,
			EventProducer: s.Producer,
		},
		Payload: TransactionCost{
			UserID:    userID,
			Source:    taskID,
			Cost:      cost,
			CreatedAt: createdAt,
		},
		PayloadSchema: s.TransactionCostSchema,
	}
}

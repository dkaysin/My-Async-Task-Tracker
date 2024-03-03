package schema_registry

import (
	"time"

	"github.com/google/uuid"
	"github.com/hamba/avro/v2"
)

type Message struct {
	Key   string
	Event Event
}

type Event struct {
	Meta
	Payload       interface{} `avro:"payload"`
	PayloadSchema avro.Schema
}

type MessageRaw struct {
	Key      string
	EventRaw EventRaw
}

// Task.Assigned

const EventNameTaskAssigned = "Task.Assigned"

func (sr *SchemaRegistry) NewEventTaskAssigned(task Task, oldUserID *string) Message {
	return Message{
		Key: task.TaskID,
		Event: Event{
			Meta: Meta{
				EventName:     EventNameTaskAssigned,
				EventID:       uuid.New().String(),
				EventVersion:  "1",
				EventProducer: sr.Producer,
			},
			Payload: TaskAssigned{
				Task:      task,
				OldUserID: oldUserID,
			},
			PayloadSchema: sr.TaskAssignedSchema,
		},
	}
}

// Task.Completed

const EventNameTaskCompleted = "Task.Completed"

func (sr *SchemaRegistry) NewEventTaskCompleted(task Task) Message {
	return Message{
		Key: task.TaskID,
		Event: Event{
			Meta: Meta{
				EventName:     EventNameTaskCompleted,
				EventID:       uuid.New().String(),
				EventVersion:  "1",
				EventProducer: sr.Producer,
			},
			Payload: TaskCompleted{
				Task: task,
			},
			PayloadSchema: sr.TaskCompletedSchema,
		},
	}
}

// Account.Created

const EventNameAccountCreated = "Account.Created"

func (sr *SchemaRegistry) NewEventAccountCreated(userID, role string) Message {
	return Message{
		Key: userID,
		Event: Event{
			Meta: Meta{
				EventName:     EventNameAccountCreated,
				EventID:       uuid.New().String(),
				EventVersion:  "1",
				EventProducer: sr.Producer,
			},
			Payload: AccountCreated{
				UserID: userID,
				Role:   role,
			},
			PayloadSchema: sr.AccountCreatedSchema,
		},
	}
}

// Account.Updated

const EventNameAccountUpdated = "Account.Updated"

func (sr *SchemaRegistry) NewEventAccountUpdated(userID, role string, active bool) Message {
	return Message{
		Key: userID,
		Event: Event{
			Meta: Meta{
				EventName:     EventNameAccountCreated,
				EventID:       uuid.New().String(),
				EventVersion:  "1",
				EventProducer: sr.Producer,
			},
			Payload: AccountUpdated{
				UserID: userID,
				Role:   role,
				Active: active,
			},
			PayloadSchema: sr.AccountCreatedSchema,
		},
	}
}

// Payment.Made

const EventNamePaymentMade = "Payment.Made"

func (sr *SchemaRegistry) NewEventPaymentMade(userID string, amount int, processedAt time.Time) Message {
	return Message{
		Key: userID,
		Event: Event{
			Meta: Meta{
				EventName:     EventNamePaymentMade,
				EventID:       uuid.New().String(),
				EventVersion:  "1",
				EventProducer: sr.Producer,
			},
			Payload: PaymentMade{
				UserID:      userID,
				Amount:      amount,
				ProcessedAt: processedAt,
			},
			PayloadSchema: sr.PaymentMadeSchema,
		},
	}
}

// Transaction.Profit

const EventNameTransactionRevenue = "Transaction.Revenue"

func (sr *SchemaRegistry) NewEventTransactionRevenue(userID *string, taskID string, revenue int, createdAt time.Time) Message {
	var key string
	if userID != nil {
		key = *userID
	}
	return Message{
		Key: key,
		Event: Event{
			Meta: Meta{
				EventName:     EventNameTransactionRevenue,
				EventID:       uuid.New().String(),
				EventVersion:  "1",
				EventProducer: sr.Producer,
			},
			Payload: TransactionRevenue{
				UserID:    userID,
				Source:    taskID,
				Revenue:   revenue,
				CreatedAt: createdAt,
			},
			PayloadSchema: sr.TransactionRevenueSchema,
		},
	}
}

const EventNameTransactionCost = "Transaction.Cost"

func (sr *SchemaRegistry) NewEventTransactionCost(userID *string, taskID string, cost int, createdAt time.Time) Message {
	var key string
	if userID != nil {
		key = *userID
	}
	return Message{
		Key: key,
		Event: Event{
			Meta: Meta{
				EventName:     EventNameTransactionCost,
				EventID:       uuid.New().String(),
				EventVersion:  "1",
				EventProducer: sr.Producer,
			},
			Payload: TransactionCost{
				UserID:    userID,
				Source:    taskID,
				Cost:      cost,
				CreatedAt: createdAt,
			},
			PayloadSchema: sr.TransactionCostSchema,
		},
	}
}

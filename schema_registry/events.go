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

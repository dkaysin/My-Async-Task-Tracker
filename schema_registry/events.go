package schema_registry

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	Key   string
	Event Event
}

type MessageRaw struct {
	Key      string
	EventRaw EventRaw
}

type Event struct {
	Meta
	Payload       interface{} `avro:"payload"`
	PayloadSchema []byte
}

type EventRaw struct {
	Meta
	Payload []byte `avro:"payload"`
}

type Meta struct {
	EventName     string `avro:"event_name"`
	EventId       string `avro:"event_id"`
	EventVersion  string `avro:"event_version"`
	EventProducer string `avro:"event_producer"`
}

type Task struct {
	TaskID      string    `json:"task_id" avro:"task_id"`
	UserID      string    `json:"user_id" avro:"user_id"`
	Description string    `json:"description" avro:"description"`
	Completed   bool      `json:"completed" avro:"completed"`
	CreatedAt   time.Time `json:"created_at" avro:"created_at"`
}

// Task.Assigned

const EventNameTaskAssigned = "Task.Assigned"

type EventPayloadTaskAssigned struct {
	Task      Task    `avro:"task"`
	OldUserID *string `avro:"old_user_id"`
}

func (sr *Schemas) NewEventTaskAssigned(task Task, oldUserID *string) Message {
	return Message{
		Key: task.TaskID,
		Event: Event{
			Meta: Meta{
				EventName:     EventNameTaskAssigned,
				EventId:       uuid.New().String(),
				EventVersion:  "1",
				EventProducer: sr.Producer,
			},
			Payload: EventPayloadTaskAssigned{
				Task:      task,
				OldUserID: oldUserID,
			},
			PayloadSchema: sr.TaskAssignedSchema,
		},
	}
}

// Task.Completed

const EventNameTaskCompleted = "Task.Completed"

type EventPayloadTaskCompleted struct {
	Task Task `avro:"task"`
}

func (sr *Schemas) NewEventTaskCompleted(task Task) Message {
	return Message{
		Key: task.TaskID,
		Event: Event{
			Meta: Meta{
				EventName:     EventNameTaskCompleted,
				EventId:       uuid.New().String(),
				EventVersion:  "1",
				EventProducer: sr.Producer,
			},
			Payload: EventPayloadTaskCompleted{
				Task: task,
			},
			PayloadSchema: sr.TaskCompletedSchema,
		},
	}
}

// Account.Created

const EventNameAccountCreated = "Account.Created"

type EventPayloadAccountCreated struct {
	UserID string `avro:"user_id"`
	Role   string `avro:"role"`
}

func (sr *Schemas) NewEventAccountCreated(userID, role string) Message {
	return Message{
		Key: userID,
		Event: Event{
			Meta: Meta{
				EventName:     EventNameAccountCreated,
				EventId:       uuid.New().String(),
				EventVersion:  "1",
				EventProducer: sr.Producer,
			},
			Payload: EventPayloadAccountCreated{
				UserID: userID,
				Role:   role,
			},
			PayloadSchema: sr.AccountCreatedSchema,
		},
	}
}

// Account.Updated

const EventNameAccountUpdated = "Account.Updated"

type EventPayloadAccountUpdated struct {
	UserID string `avro:"user_id"`
	Active bool   `avro:"active"`
	Role   string `avro:"role"`
}

func (sr *Schemas) NewEventAccountUpdated(userID, role string, active bool) Message {
	return Message{
		Key: userID,
		Event: Event{
			Meta: Meta{
				EventName:     EventNameAccountCreated,
				EventId:       uuid.New().String(),
				EventVersion:  "1",
				EventProducer: sr.Producer,
			},
			Payload: EventPayloadAccountUpdated{
				UserID: userID,
				Role:   role,
				Active: active,
			},
			PayloadSchema: sr.AccountCreatedSchema,
		},
	}
}

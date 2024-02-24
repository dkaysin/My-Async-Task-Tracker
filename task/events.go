package task

type Event struct {
	Key   string
	Value interface{}
}

// Task.Assigned

const EventKeyTaskAssigned = "Task.Assigned"

func NewEventTaskAssigned(task Task) Event {
	return Event{
		Key:   EventKeyTaskAssigned,
		Value: task,
	}
}

// Task.Completed

const EventKeyTaskCompleted = "Task.Completed"

func NewEventTaskCompleted(task Task) Event {
	return Event{
		Key:   EventKeyTaskCompleted,
		Value: task,
	}
}

// Account.Created

const EventKeyAccountCreated = "Account.Created"

type EventValueAccountCreateAccount struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

// Account.Updated

const EventKeyAccountUpdated = "Account.Updated"

type EventValueAccountUpdated struct {
	UserID string `json:"user_id"`
	Active bool   `json:"active"`
	Role   string `json:"role"`
}

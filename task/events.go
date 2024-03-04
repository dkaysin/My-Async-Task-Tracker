package task

type Event struct {
	Key   string
	Value interface{}
}

// Task.Assigned

const EventKeyTaskAssigned = "Task.Assigned"

func NewEventTaskAssigned(task Task, oldUserID *string) Event {
	return Event{
		Key: EventKeyTaskAssigned,
		Value: EventValueTaskAssigned{
			Task:      task,
			OldUserID: oldUserID,
		},
	}
}

type EventValueTaskAssigned struct {
	Task
	OldUserID *string `json:"old_user_id"`
}

// Task.Completed

const EventKeyTaskCompleted = "Task.Completed"

func NewEventTaskCompleted(value Task) Event {
	return Event{
		Key:   EventKeyTaskCompleted,
		Value: value,
	}
}

// Account.Created

const EventKeyAccountCreated = "Account.Created"

type EventValueAccountCreated struct {
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

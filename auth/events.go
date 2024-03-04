package auth

type Event struct {
	Key   string
	Value interface{}
}

// Account.Created

const EventKeyAccountCreated = "Account.Created"

type EventValueAccountCreateAccount struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
}

func NewEventAccountCreated(userID, role string) Event {
	return Event{
		Key: EventKeyAccountCreated,
		Value: EventValueAccountCreateAccount{
			UserID: userID,
			Role:   role,
		},
	}
}

// Account.Updated

const EventKeyAccountUpdated = "Account.Updated"

type EventValueAccountUpdated struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	Active bool   `json:"active"`
}

func NewEventAccountUpdated(userID, role string, active bool) Event {
	return Event{
		Key: EventKeyAccountCreated,
		Value: EventValueAccountUpdated{
			UserID: userID,
			Role:   role,
			Active: active,
		},
	}
}

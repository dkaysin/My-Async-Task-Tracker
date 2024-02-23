package auth

import "errors"

// kafka
const (
	KafkaConsumerGroupID = "my-consumer-group-id"

	KafkaTopicIDA = "topic-A"
	KafkaTopicIDB = "topic-B"
)

// events
const (
	Event1 = "event_1"
	Event2 = "event_2"
)

type AddUserReq struct {
	UserId string `json:"user_id" validate:"required"`
}


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

// errors
var ErrPayloadValidationFailed = errors.New("payload validation failed")

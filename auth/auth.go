package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

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

type JwtCustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// errors
var ErrPayloadValidationFailed = errors.New("payload validation failed")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidJwtClaimsFormat = errors.New("invalid jwt claims format")
var ErrInsufficientPrivileges = errors.New("insufficient privileges")
var ErrTokenNotFound = errors.New("token not found in request context")

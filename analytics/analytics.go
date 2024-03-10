package analytics

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

const (
	RoleDeveloper  = "developer"
	RoleAdmin      = "admin"
	RoleManager    = "manager"
	RoleAccountant = "accountant"
)

// kafka
const KafkaConsumerGroupID = "consumer-group-analytics"
const ProducerName = "Analytics"

// cross-package types

type JwtCustomClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// errors
var ErrPayloadValidationFailed = errors.New("payload validation failed")
var ErrInvalidJwtClaimsFormat = errors.New("invalid jwt claims format")
var ErrInsufficientPrivileges = errors.New("insufficient privileges")
var ErrTokenNotFound = errors.New("token not found in request context")
var ErrTaskNotFound = errors.New("task not found")
var ErrNoDevelopersAvailable = errors.New("no developers available")
var ErrUnknownUser = errors.New("unknown user")
var ErrMessageHeaderNotFound = errors.New("required kafka message header not found")
var ErrUnknownEventVersion = errors.New("unknown event version")

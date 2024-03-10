package task

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	RoleDeveloper  = "developer"
	RoleAdmin      = "admin"
	RoleManager    = "manager"
	RoleAccountant = "accountant"
)

// kafka
const KafkaConsumerGroupID = "consumer-group-task"
const ProducerName = "Task"

// cross-package types

type JwtCustomClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type Task struct {
	TaskID      string    `json:"task_id"`
	UserID      string    `json:"user_id"`
	Description string    `json:"description"`
	JiraID      string    `json:"jira_id"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
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

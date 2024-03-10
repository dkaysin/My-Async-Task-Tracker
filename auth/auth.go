package auth

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
const KafkaConsumerGroupID = "consumer-group-auth"
const ProducerName = "Auth"

// cross-package types

type JwtCustomClaims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// errors
var ErrPayloadValidationFailed = errors.New("payload validation failed")
var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrInvalidJwtClaimsFormat = errors.New("invalid jwt claims format")
var ErrInsufficientPrivileges = errors.New("insufficient privileges")
var ErrTokenNotFound = errors.New("token not found in request context")
var ErrAccountNotFound = errors.New("account not found")

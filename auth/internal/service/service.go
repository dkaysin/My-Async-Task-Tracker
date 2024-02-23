package service

import (
	database "async_course/auth/internal/database"
	writer "async_course/auth/internal/event_writer"

	"github.com/spf13/viper"
)

type Service struct {
	config     *viper.Viper
	db         *database.Database
	ew         *writer.EventWriter
	signingKey string
}

func NewService(config *viper.Viper, db *database.Database, ew *writer.EventWriter, signingKey string) *Service {
	return &Service{
		config:     config,
		db:         db,
		ew:         ew,
		signingKey: signingKey,
	}
}

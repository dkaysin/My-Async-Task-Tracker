package service

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v5"
	"github.com/segmentio/kafka-go"
)

func (s *Service) InsertIntoDLQ(m kafka.Message, err error) error {
	headersBytes, err := json.Marshal(m.Headers)
	if err != nil {
		return err
	}
	return s.db.ExecuteTx(context.Background(), func(tx pgx.Tx) error {
		q := `INSERT INTO accounting_dlq (message_headers, message_key, message_value, error) VALUES ($1, $2, $3, $4)`
		_, err := tx.Exec(context.Background(), q, string(headersBytes), m.Key, m.Value, err.Error())
		return err
	})
}

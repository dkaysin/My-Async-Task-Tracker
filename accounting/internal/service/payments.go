package service

import "time"

func (s *Service) processPayment(userID string, amount int) (time.Time, error) {
	// call external service to process payment
	return time.Now().UTC(), nil
}

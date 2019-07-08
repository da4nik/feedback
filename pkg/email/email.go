package email

import (
	"fmt"
)

// Email - email entity
type Email struct {
}

// NewEmail - returns new email instance
func NewEmail() (Email, error) {
	return Email{}, nil
}

// Send - sends email
func (e Email) Send(email, subject, body string) error {
	return fmt.Errorf("not yet implemented")
}

package email

import (
	"fmt"
)

// Email - email entity
type Email struct {
}

// NewEmail - returns new email instance
func NewEmail(targetEmail, mandrillKey string) (Email, error) {
	return Email{}, nil
}

// Send - sends email
func (e Email) Send(email, message string) error {
	return fmt.Errorf("not yet implemented")
}

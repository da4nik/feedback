package email

import "fmt"

// Email - email entity
type Email struct{}

// NewEmail - returns new email instance
func NewEmail() Email {
	return Email{}
}

// Send - sends email
func (e Email) Send(email, message string) error {
	return fmt.Errorf("not yet implemented")
}

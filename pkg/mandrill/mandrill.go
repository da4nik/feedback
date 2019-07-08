package mandrill

import (
	"github.com/da4nik/feedback/internal/log"
	"github.com/mattbaird/gochimp"
)

// Mandrill - email entity via mandrill
type Mandrill struct {
	mandrillAPI *gochimp.MandrillAPI
	targetEmail string
}

// NewMandrill - returns new mandrill email instance
func NewMandrill(mandrillKey string) (Mandrill, error) {
	mandrillAPI, err := gochimp.NewMandrill(mandrillKey)
	if err != nil {
		return Mandrill{}, err
	}

	return Mandrill{
		mandrillAPI: mandrillAPI,
	}, nil
}

// Send - sends email
func (m Mandrill) Send(email, subject, body string) error {
	message := gochimp.Message{
		Text:      body,
		Subject:   subject,
		FromEmail: "feedback@captureproof.com",
		FromName:  "Marketing Site",
		To: []gochimp.Recipient{
			gochimp.Recipient{
				Email: email,
				Name:  "Sales guys",
				Type:  "",
			},
		},
	}

	log.Debugf("Sending email: %s <- [%s] %s", email, subject, body)

	if _, err := m.mandrillAPI.MessageSend(message, false); err != nil {
		log.Errorf("error sending mandrill email: %s", err.Error())
		return err
	}
	return nil
}

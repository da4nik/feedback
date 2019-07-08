package config

import (
	"os"
	"strconv"
)

// Config application cofig entity
type Config struct {
	MandrillKey string
	TargetEmail string
	Port        int
	TwilioSID   string
	TwilioKey   string
	TwilioPhone string
}

// LoadConfig - load config
func LoadConfig() Config {
	serverPort, err := strconv.Atoi(os.Getenv("FEEDBACK_PORT"))
	if err != nil {
		serverPort = 9000
	}

	targetEmail := os.Getenv("FEEDBACK_TARGET_EMAIL")
	if len(targetEmail) == 0 {
		targetEmail = "sales@captureproof.com"
	}

	twilioPhone := os.Getenv("TWILIO_PHONE")
	if len(twilioPhone) == 0 {
		twilioPhone = "+14157671697"
	}

	return Config{
		TargetEmail: targetEmail,
		MandrillKey: os.Getenv("MANDRILL_KEY"),
		Port:        serverPort,
		TwilioSID:   os.Getenv("TWILIO_SID"),
		TwilioKey:   os.Getenv("TWILIO_KEY"),
		TwilioPhone: twilioPhone,
	}
}

package text

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// Text entity
type Text struct {
	authToken  string
	accountSid string
	fromNumber string
	baseURL    string
}

// NewText returns new Text instance
func NewText(fromNumber, accountSid, authToken string) Text {
	return Text{
		accountSid: accountSid,
		authToken:  authToken,
		fromNumber: fromNumber,
		baseURL:    "https://api.twilio.com/2010-04-01",
	}
}

// Send sends text to phone number
func (t Text) Send(toNumber, text string) error {
	if !isValidNumber(toNumber) {
		return fmt.Errorf("phone number is not valid: %s", toNumber)
	}

	apiURL := t.url("/Messages.json")

	v := url.Values{}
	v.Set("To", toNumber)
	v.Set("From", t.fromNumber)
	v.Set("Body", text)

	rb := *strings.NewReader(v.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", apiURL, &rb)
	req.SetBasicAuth(t.accountSid, t.authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode == 201 {
		return nil
	}

	// var data map[string]interface{}
	// bodyBytes, _ := ioutil.ReadAll(resp.Body)
	// json.Unmarshal(bodyBytes, &data)

	return fmt.Errorf("twilio api http code: %s", resp.Status)
}

// isValidNumber verifies phone number with twillio api
func isValidNumber(phone string) bool {
	matched, _ := regexp.MatchString(`^\+?[1-9]\d{1,14}$`, phone)
	return matched
}

func (t Text) url(path string) string {
	return fmt.Sprintf("%s/Accounts/%s%s", t.baseURL, t.accountSid, path)
}

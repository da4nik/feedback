package handlers

import (
	"fmt"
	"net/http"

	"github.com/da4nik/feedback/internal/log"
)

func (h Handlers) knowADoctor(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("name")
	location := r.PostFormValue("location")
	specialty := r.PostFormValue("specialty")
	phone := r.PostFormValue("phone")

	text := fmt.Sprintf(
		"A patient knows a doctor %s (%s), who is working at %s. "+
			"Doctor's phone number is %s",
		name,
		specialty,
		location,
		phone)

	log.Debugf("Feedback email, doctor %q (%s)", name, phone)

	err := h.email.Send(
		"sales@captureproof.com",
		"[marketing] Patient knows a doctor",
		text)
	if err != nil {
		log.Errorf("Unable to send email: %s", err.Error())
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	w.WriteHeader(http.StatusOK)
}

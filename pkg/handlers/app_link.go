package handlers

import (
	"net/http"
	"strings"

	"github.com/da4nik/feedback/internal/log"
)

const iOSAppLinkMessage = "Captureproof iOS app https://itunes.apple.com/app/captureproof/id987584313"
const androidAppLinkMessage = "Captureproof Android app https://play.google.com/store/apps/details?id=com.captureproof&hl=en"

func (h Handlers) appLink(w http.ResponseWriter, r *http.Request) {
	phone := r.PostFormValue("phone")
	phone = strings.Trim(phone, " \n")

	if len(phone) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Debugf("Sending app link to: %s", phone)
	if err := h.sendTextLink(phone, iOSAppLinkMessage, w); err != nil {
		log.Errorf("unable to send ios link: %s", err.Error())
		return
	}

	log.Debugf("Sending android app link to: %s", phone)
	if err := h.sendTextLink(phone, androidAppLinkMessage, w); err != nil {
		log.Errorf("unable to send android link: %s", err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h Handlers) sendTextLink(phone, link string, w http.ResponseWriter) error {
	if err := h.text.Send(phone, link); err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return err
	}
	return nil
}

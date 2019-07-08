package handlers

import (
	"net/http"
	"strings"
)

type text interface {
	Send(string, string) error
}

type email interface {
	Send(email, subject, body string) error
}

type handlers = map[string]http.HandlerFunc

// Handlers represents handlers entity
type Handlers struct {
	text     text
	email    email
	handlers handlers
}

// NewHandlers - return new instance of Handlers
func NewHandlers(text text, email email) Handlers {
	h := Handlers{
		text:  text,
		email: email,
	}

	handlerMap := handlers{
		"/know-a-doctor": h.knowADoctor,
		"/app-link":      h.appLink,
	}

	h.handlers = handlerMap

	return h
}

// ServerHTTP servers http request
func (h Handlers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Middlewares
	fun :=
		logger(
			h.cors(
				h.onlyAllowedMethod("POST",
					h.processEndpoints)))
	fun(w, r)
}

func (h Handlers) processEndpoints(w http.ResponseWriter, r *http.Request) {
	if fun, exists := h.handlers[r.URL.Path]; exists {
		fun(w, r)
		return
	}

	w.WriteHeader(http.StatusNotFound)
}

func (h Handlers) onlyAllowedMethod(
	method string,
	next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		next(w, r)
	}
}

func (h Handlers) cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		origin := strings.Trim(r.Header.Get("Origin"), " ")
		if !h.isValidHost(origin) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, Origin")
		w.Header().Set("Access-Control-Allow-Methods", "POST")

		next(w, r)
	}
}

func (h Handlers) isValidHost(host string) (result bool) {
	result = strings.Contains(host, "captureproof.com") ||
		strings.Contains(host, "cp2.div-art.com.ua") ||
		strings.Contains(host, "localhost")
	return
}

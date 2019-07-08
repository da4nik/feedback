package handlers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/da4nik/feedback/internal/log"
)

type responseWriterProxy struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (o *responseWriterProxy) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *responseWriterProxy) WriteHeader(code int) {
	if o.wroteHeader {
		log.Errorf("Header double writing not allowed")
		return
	}

	o.ResponseWriter.WriteHeader(code)

	o.wroteHeader = true
	o.status = code
}

func logger(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rwp := &responseWriterProxy{ResponseWriter: w}

		h.ServeHTTP(rwp, r)

		addr := r.RemoteAddr
		if i := strings.LastIndex(addr, ":"); i != -1 {
			addr = addr[:i]
		}

		log.Infof("%q %d %d %q %q %s",
			fmt.Sprintf("%s %s %s", r.Method, r.URL, r.Proto),
			rwp.status,
			rwp.written,
			r.Referer(),
			r.UserAgent(),
			time.Since(start))
	})
}

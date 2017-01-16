package middleware

import (
	"log"
	"net/http"
)

func LoggedHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s requested %s", r.RemoteAddr, r.URL)
		h.ServeHTTP(w, r)
	})
}

func HeaderedHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s requested %s", r.RemoteAddr, r.URL)
		w.Header().Set("powered-by", "go-ocelot")
		if r.TLS != nil {
			r.Header.Add("x-forwarded-proto", "https")
		} else {
			r.Header.Add("x-forwarded-proto", "http")
		}
		h.ServeHTTP(w, r)
	})
}
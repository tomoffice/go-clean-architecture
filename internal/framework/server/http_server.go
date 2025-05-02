package server

import (
	"log"
	"net/http"
)

func StartHTTPServer(addr string, handler http.Handler) {
	log.Printf("Server starting at %s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("Server errordefs: %v", err)
	}
}

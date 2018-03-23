package webserver

import (
	"io"
	"net/http"
)

func webHealth(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Running")
}

// StartWebServer exports API and healh check information
func StartWebServer() {
	// Run webserver in backend
	go func() {
		http.HandleFunc("/_ping", webHealth)
		http.ListenAndServe(":8080", nil)
	}()
}

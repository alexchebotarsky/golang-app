package handler

import (
	"io"
	"log"
	"net/http"
)

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := io.WriteString(w, http.StatusText(http.StatusOK)); err != nil {
		log.Printf("%s: %v", logMsgWriteResponse, err)
	}
}

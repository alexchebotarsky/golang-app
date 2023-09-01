package handler

import "log"

func handleWritingErr(err error) {
	if err != nil {
		log.Printf("Error writing to http.ResponseWriter: %v", err)
	}
}

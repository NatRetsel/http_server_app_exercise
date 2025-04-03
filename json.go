package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	/*
		1.) Logs error if error is not nil
		2.) Capture server errors
		3.) Respond with JSON
			{
				error: "${error message}"
			}
	*/
	if err != nil {
		log.Println(err)
	}
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, errorResponse{
		Error: msg,
	})

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	/*
		1.) Marshal the payload into JSON format
		2.) Send code and response via the http.ResponseWriter
	*/
	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(payload)
	if err != nil {
		log.Printf("error marshalling json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(jsonResp)
}

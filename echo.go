package main

import (
	"encoding/json"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type EchoRequest struct {
	Message string `json:"Message"`
}

func (p EchoRequest) String() string {
	return "echo: Message: " + p.Message
}

type EchoResponse struct {
	Date    string `json:"Date"`
	Message string `json:"Message"`
}

func Echo(w http.ResponseWriter, r *http.Request) {
	var echoRequest EchoRequest
	var echoResponse EchoResponse

	//logger.WithFields(logger.Fields{
	//	"application": "restsvc",
	//	"service": "echo",
	//}).Info("echo service")
	logger.Info("echo service called")

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(&echoRequest)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Debug(echoRequest)

	// Return the current date and the original message
	echoResponse.Date = time.Now().UTC().Format(time.RFC3339)
	echoResponse.Message = echoRequest.Message

	// Write the response
	j, err := json.Marshal(echoResponse)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

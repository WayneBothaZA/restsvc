package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/nullseed/logruseq"
	logger "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const APP_VERSION = "0.0.2 (debug)"
const URL = ":8080"

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "print the version number.")
var debugFlag *bool = flag.Bool("d", false, "enable debug logging")

func init() {
	// Log to seq instance
	logger.AddHook(logruseq.NewSeqHook("http://localhost:5341"))
	logger.SetOutput(os.Stdout)

	// Setup signal handler to cleanup properly
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Start signal handler thread
	go signalHandler(signals)
}

func shutdown() {
	logger.Info("Stopping REST service on " + URL)
}

func signalHandler(signals <-chan os.Signal) {
	sig := <-signals
	logger.Error("signal: ", sig)
	shutdown()
	os.Exit(1)
}

func main() {
	flag.Parse() // Scan the arguments list

	// Display application version
	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
		os.Exit(0)
	}

	if *debugFlag {
		logger.SetLevel(logger.DebugLevel)
		logger.SetReportCaller(true)
	}
	logger.Info("Starting REST service on " + URL)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/echo", Echo).Methods(http.MethodPost)
	loggedRouter := handlers.LoggingHandler(logger.StandardLogger().Writer(), router)
	logger.Fatal(http.ListenAndServe(URL, loggedRouter))
}

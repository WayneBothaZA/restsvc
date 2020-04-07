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

const APP_VERSION = "1.0.2"
const DEFAULT_URL = ":80"

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "print the version number.")
var debugFlag *bool = flag.Bool("d", false, "enable debug logging")
var dockerMode *bool = flag.Bool("docker", false, "enable docker mode")

// URL (port) on which to start our listener
var url string

func init() {
	// Log to seq instance
	seq, found := os.LookupEnv("SEQ")
	if found {
		logger.AddHook(logruseq.NewSeqHook(seq))
	}
	logger.SetOutput(os.Stdout)

	// set loglevel
	switch os.Getenv("LOGLEVEL") {
	case "DEBUG":
		logger.SetLevel(logger.DebugLevel)
	case "WARN":
		logger.SetLevel(logger.WarnLevel)
	default:
		logger.SetLevel(logger.InfoLevel)
	}

	// configure
	url, found = os.LookupEnv("URL")
	if !found {
		url = DEFAULT_URL
	}

	// Setup signal handler to cleanup properly
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Start signal handler thread
	go signalHandler(signals)
}

func shutdown() {
	logger.Info("Stopping REST service on " + url)
}

func signalHandler(signals <-chan os.Signal) {
	sig := <-signals
	logger.Error("signal: ", sig)
	shutdown()
	os.Exit(1)
}

func main() {
	// scan uarguments
	flag.Parse()

	// display application version
	if *versionFlag {
		fmt.Println("Version:", APP_VERSION)
		os.Exit(0)
	}

	logger.Info("Starting REST service on " + url)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/echo", Echo).Methods(http.MethodPost)
	loggedRouter := handlers.LoggingHandler(logger.StandardLogger().Writer(), router)
	logger.Fatal(http.ListenAndServe(url, loggedRouter))
}

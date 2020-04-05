package main

import (
	"flag"
	"fmt"
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
	logger.AddHook(logruseq.NewSeqHook("http://localhost:5341"))
	// Or optionally use the hook with an API key:
	// log.AddHook(logruseq.NewSeqHook("http://localhost:5341",
	// 	logruseq.OptionAPIKey("N1ncujiT5pYGD6m4CF0")))

	// Log as JSON instead of the default ASCII formatter.
	//logger.SetFormatter(&logger.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)
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

	// Setup signal handler to cleanup properly
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	// Start signal handler thread
	go signalHandler(signals)

	logger.Info("Starting REST service on " + URL)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/echo", Echo).Methods(http.MethodPost)
	logger.Fatal(http.ListenAndServe(URL, router))
}

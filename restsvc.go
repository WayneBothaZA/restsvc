package main 

import (
	"os"
	"os/signal"
	"syscall"
	"flag"
	"fmt"
	"time"
	"net/http"
	"encoding/json"

	"github.com/gorilla/mux"
	"github.com/nullseed/logruseq"
	logger "github.com/sirupsen/logrus"
)

const APP_VERSION = "0.0.1 (debug)"
const URL = ":8080"

// The flag package provides a default help printer via -h switch
var versionFlag *bool = flag.Bool("v", false, "Print the version number.")
//var versionFlag *bool = flag.Bool("v", false, "Print the version number.")

type PingRequest struct {
        Message string `json:"Message"`
}

func (p PingRequest) String() string {
    return "Ping: Message: " + p.Message
}

type PingResponse struct {
		Date string `json:"Date"`
        Message string `json:"Message"`
}

func ping(w http.ResponseWriter, r *http.Request) {
	var preq PingRequest
	var presp PingResponse
	
	//logger.WithFields(logger.Fields{
	//	"application": "restsvc",
	//	"service": "ping",
	//}).Info("ping service")
	logger.Info("ping service called")

    // Try to decode the request body into the struct. If there is an error,
    // respond to the client with the error message and a 400 status code.
    err := json.NewDecoder(r.Body).Decode(&preq)
    if err != nil {
        logger.Error(err.Error())
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    logger.Info(preq)

	presp.Date = time.Now().UTC().Format(time.RFC3339)
	presp.Message = preq.Message
	
	// Write the response
	w.Header().Set("Content-Type", "application/json")
    
    j, err := json.Marshal(presp)
    if err != nil {
    	logger.Error(err.Error())
    	http.Error(w, err.Error(), http.StatusInternalServerError)
    	return
    }

	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func init() {
	logger.AddHook(logruseq.NewSeqHook("http://localhost:5341"))
	logger.SetReportCaller(true)
	// Or optionally use the hook with an API key:
	// log.AddHook(logruseq.NewSeqHook("http://localhost:5341",
	// 	logruseq.OptionAPIKey("N1ncujiT5pYGD6m4CF0")))

	// Log as JSON instead of the default ASCII formatter.
	//logger.SetFormatter(&log.JSONFormatter{})
	
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)
	
	// Only log the warning severity or above.
	//logger.SetLevel(log.WarnLevel)
}

func shutdown() {
	logger.Info("Done")
}

func signalHandler(signals <-chan os.Signal) {
    	sig := <- signals
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
   
	// Setup signal handler to cleanup properly
	signals := 	make(chan os.Signal, 1)
    signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
    
    // Start signal handler thread
    go signalHandler(signals)
   
	//logger.WithFields(logger.Fields{
	//	"application": "restsvc",
	//}).Info("Starting REST service" + URL)
	logger.Info("Starting REST service...")
	
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/ping", ping).Methods(http.MethodPost)
	logger.Fatal(http.ListenAndServe(URL, router))
}

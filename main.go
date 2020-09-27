package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	serviceUri      = "https://api.us-west-1.saucelabs.com/v1/magnificent/"
	defaultInterval = 10
)

var (
	serviceFailures     = int64(0)
	serviceOks          = int64(0)
	totalCount          = int64(0)
	serviceUnresponsive = int64(0)

	magnificent = NewMagnificentClient(serviceUri)
)

// MagnificentClient the client struct to put the methods on
type MagnificentClient struct {
	baseURL  string
	client   *http.Client
	mustStop bool
}

// StatusResponse the response struct to export in json
type StatusResponse struct {
	ServiceFailures     int64   `json:"service_failures"`
	ServiceOks          int64   `json:"service_oks"`
	TotalCount          int64   `json:"total_count"`
	ServiceUnresponsive int64   `json:"service_unresponsive"`
	Availability        float64 `json:"availability"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", getStatus).Methods(http.MethodGet)
	r.HandleFunc("/callit", callIt).Methods(http.MethodGet)
	r.HandleFunc("/muststop", mustStop).Methods(http.MethodGet)
	interval := defaultInterval
	if len(os.Args[1:]) > 0 {
		interval, _ = strconv.Atoi(os.Args[1])
	}
	log.Println("Starting monitoring magnificent every", interval, "seconds")
	go runsMagnificent(magnificent, interval)
	log.Fatal(http.ListenAndServe(":8080", r))

}

// getStatus writes the status to the response in json format
func getStatus(w http.ResponseWriter, r *http.Request) {
	response := StatusResponse{
		ServiceFailures:     serviceFailures,
		ServiceOks:          serviceOks,
		TotalCount:          totalCount,
		ServiceUnresponsive: serviceUnresponsive,
		Availability:        float64(serviceOks) / float64(totalCount),
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

// callIt calls magnificent and writes the return to the reponse
func callIt(w http.ResponseWriter, r *http.Request) {
	status, err := magnificent.callMagnificient()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error in calling service"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(status))
}

// mustStop sets the mustStop boolean to true, so that the monitoring stops
func mustStop(w http.ResponseWriter, r *http.Request) {
	magnificent.mustStop = true
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Stop boolean set!"))
}

// NewMagnificentClient creates a new instance of the mag client
func NewMagnificentClient(url string) *MagnificentClient {
	return &MagnificentClient{
		baseURL:  url,
		client:   http.DefaultClient,
		mustStop: false,
	}
}

// callMagnificient calls the service and updates the counters
func (m *MagnificentClient) callMagnificient() (string, error) {
	req, err := http.NewRequest("GET", m.baseURL, nil)
	if err != nil {
		log.Println("Error creating the request ", err)
		return "", err
	}
	resp, err := m.client.Do(req)
	totalCount++
	if err != nil {
		log.Println("Error sending the request. ", err)
		serviceUnresponsive++
		return "Service unresponsive", err
	}

	if resp.StatusCode == 200 {
		serviceOks++
		return "OK", nil
	}
	if resp.StatusCode == 500 {
		serviceFailures++
		return "ERROR", nil
	}
	return "Magnificent return not understood", nil
}

// runsMagnificent runs the call to magnificent as long as the stop boolean is false
func runsMagnificent(client *MagnificentClient, interval int) {
	for !client.mustStop {
		client.callMagnificient()
		time.Sleep(time.Duration(interval) * time.Second)
	}
}

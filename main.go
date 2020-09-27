package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

const (
	serviceUri = "https://api.us-west-1.saucelabs.com/v1/magnificent/"
)

var (
	serviceFailures = int64(0)
	serviceOks      = int64(0)
	count           = int64(0)
	magnificent = NewMagnificentClient(serviceUri)
)

// MagnificentClient the client struct to put the methods on
type MagnificentClient struct {
	baseURL string
	client  *http.Client
}

// StatusResponse the response struct to export in json
type StatusResponse struct {
	ServiceFailures int64 `json:"service_failures"`
	ServiceOks      int64 `json:"service_oks"`
	TotalCount      int64 `json:"total_count"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", getStatus).Methods(http.MethodGet)
	r.HandleFunc("/callit", callIt).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", r))

}

// getStatus writes the status to the response in json format
func getStatus(w http.ResponseWriter, r *http.Request) {
	response := StatusResponse{
		ServiceFailures: serviceFailures,
		ServiceOks:      serviceOks,
		TotalCount:      count,
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		panic(err)
	}
}

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

func NewMagnificentClient(url string) *MagnificentClient {
	return &MagnificentClient{
		baseURL: url,
		client:  http.DefaultClient,
	}
}

func (m *MagnificentClient) callMagnificient() (string, error) {
	req, err := http.NewRequest("GET", m.baseURL, nil)
	if err != nil {
		log.Println("Error creating the request %s", err)
		return "", err
	}
	resp, err := m.client.Do(req)

	if err != nil {
		log.Println("Error sending the request. %s", err)
		return "", err
	}
	count++

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

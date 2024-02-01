package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Struct of Request to represent the incoming JSON data
type Request struct {
	Ev     string `json:"ev"`
	Et     string `json:"et"`
	ID     string `json:"id"`
	UID    string `json:"uid"`
	MID    string `json:"mid"`
	T      string `json:"t"`
	P      string `json:"p"`
	L      string `json:"l"`
	SC     string `json:"sc"`
	ATRK1  string `json:"atrk1"`
	ATRV1  string `json:"atrv1"`
	ATRT1  string `json:"atrt1"`
	ATRK2  string `json:"atrk2"`
	ATRV2  string `json:"atrv2"`
	ATRT2  string `json:"atrt2"`
	UATRK1 string `json:"uatrk1"`
	UATRV1 string `json:"uatrv1"`
	UATRT1 string `json:"uatrt1"`
	UATRK2 string `json:"uatrk2"`
	UATRV2 string `json:"uatrv2"`
	UATRT2 string `json:"uatrt2"`
	UATRK3 string `json:"uatrk3"`
	UATRV3 string `json:"uatrv3"`
	UATRT3 string `json:"uatrt3"`
}

// ConvertedRequest struct to represent the converted JSON data
type ConvertedRequest struct {
	Event           string               `json:"event"`
	EventType       string               `json:"event_type"`
	AppID           string               `json:"app_id"`
	UserID          string               `json:"user_id"`
	MessageID       string               `json:"message_id"`
	PageTitle       string               `json:"page_title"`
	PageURL         string               `json:"page_url"`
	BrowserLanguage string               `json:"browser_language"`
	ScreenSize      string               `json:"screen_size"`
	Attributes      map[string]Attribute `json:"attributes"` // map of attribute
	Traits          map[string]Trait     `json:"traits"`     // map of trait
}

type Attribute struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Trait struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

func main() {
	// Create a channel to send requests to the worker
	requestChannel := make(chan Request)

	// Start the worker
	go worker(requestChannel)

	// Setup HTTP - handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//POST Method
		if r.Method == "POST" {
			// Parse the JSON request
			var req Request
			err := json.NewDecoder(r.Body).Decode(&req)
			if err != nil {
				http.Error(w, "Error parsing JSON", http.StatusBadRequest)
				return
			}

			// Send the request to the worker via the channel
			requestChannel <- req

			// Respond to the client
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Request received successfully"))
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	})

	// Start the HTTP server on PORT 8080
	fmt.Println("Server listening on :8080")
	http.ListenAndServe(":8080", nil)
}

func worker(channel <-chan Request) {
	for {
		// Wait for a request from the channel
		req := <-channel

		// Launch a new Goroutine to process each request
		go processRequest(req)

	}
}

// This function is executed in a separate Goroutine for each incoming request
func processRequest(req Request) {
	// fmt.Println("Processing request:", req)

	// Convert the request to the desired format
	convertedReq := convertRequest(req)

	// Send the converted request to the webhook
	err := sendToWebhook(convertedReq)
	if err != nil {
		fmt.Println("Error sending to webhook:", err)
	}
}

func convertRequest(req Request) ConvertedRequest {
	convertedReq := ConvertedRequest{
		Event:           req.Ev,
		EventType:       req.Et,
		AppID:           req.ID,
		UserID:          req.UID,
		MessageID:       req.MID,
		PageTitle:       req.T,
		PageURL:         req.P,
		BrowserLanguage: req.L,
		ScreenSize:      req.SC,
		Attributes: map[string]Attribute{
			req.ATRK1: {Value: req.ATRV1, Type: req.ATRT1},
			req.ATRK2: {Value: req.ATRV2, Type: req.ATRT2},
		},
		Traits: map[string]Trait{
			req.UATRK1: {Value: req.UATRV1, Type: req.UATRT1},
			req.UATRK2: {Value: req.UATRV2, Type: req.UATRT2},
			req.UATRK3: {Value: req.UATRV3, Type: req.UATRT3},
		},
	}

	return convertedReq
}

func sendToWebhook(convertedReq ConvertedRequest) error {
	// Convert the ConvertedRequest to JSON
	jsonData, err := json.Marshal(convertedReq)
	if err != nil {
		return err
	}

	// Send the JSON data to the webhook
	resp, err := http.Post("https://webhook.site/10c63412-e94e-446b-96ea-9cc936e4eca1", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	// close response after function executes
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Webhook request failed with status: %s", resp.Status)
	}

	fmt.Println("Sent to webhook successfully")

	return nil
}

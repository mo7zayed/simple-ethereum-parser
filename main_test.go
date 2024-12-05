package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mo7zayed/trustwallet/parser"
	"github.com/mo7zayed/trustwallet/parser/ethereum"
)

func TestCurrentBlockEndpoint(t *testing.T) {
	ethereumParser := ethereum.NewEthereumParser()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentBlock := ethereumParser.GetCurrentBlock()
		response := map[string]int{"current_block": currentBlock}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}

	defer resp.Body.Close()

	var result map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if _, exists := result["current_block"]; !exists {
		t.Error("Expected 'current_block' key in response")
	}
}

func TestSubscribeEndpoint(t *testing.T) {
	ethereumParser := ethereum.NewEthereumParser()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var request struct {
			Address string `json:"address"`
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		ethereumParser.Subscribe(request.Address)
		response := map[string]string{"message": "Address subscribed successfully"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Valid request
	requestBody := []byte(`{"address": "0x123456789abcdef"}`)
	resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("Failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result["message"] != "Address subscribed successfully" {
		t.Errorf("Expected success message, got: %s", result["message"])
	}
}

func TestTransactionsEndpoint(t *testing.T) {
	ethereumParser := ethereum.NewEthereumParser()
	ethereumParser.Subscribe("0x123456789abcdef")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Address query parameter is required", http.StatusBadRequest)
			return
		}

		transactions := ethereumParser.GetTransactions(address)
		response := map[string][]parser.Transaction{"transactions": transactions}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Valid request
	resp, err := http.Get(server.URL + "?address=0x123456789abcdef")
	if err != nil {
		t.Fatalf("Failed to send GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	var result map[string][]parser.Transaction
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if _, exists := result["transactions"]; !exists {
		t.Error("Expected 'transactions' key in response")
	}
}

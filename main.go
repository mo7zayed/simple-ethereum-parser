package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mo7zayed/trustwallet/parser"
	"github.com/mo7zayed/trustwallet/parser/ethereum"
)

func main() {
	ethereumParser := ethereum.NewEthereumParser()

	http.HandleFunc("/current_block", func(w http.ResponseWriter, r *http.Request) {
		currentBlock := ethereumParser.GetCurrentBlock()
		response := map[string]int{"current_block": currentBlock}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
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

		if request.Address == "" {
			http.Error(w, "Address is required", http.StatusBadRequest)
			return
		}

		ethereumParser.Subscribe(request.Address)

		response := map[string]string{"message": "Address subscribed successfully"}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		address := r.URL.Query().Get("address")
		if address == "" {
			http.Error(w, "Address query parameter is required", http.StatusBadRequest)
			return
		}

		transactions := ethereumParser.GetTransactions(address)

		response := map[string][]parser.Transaction{"transactions": transactions}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	})

	log.Println("Server is running on port 3000...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

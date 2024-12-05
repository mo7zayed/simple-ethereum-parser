package ethereum

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/mo7zayed/trustwallet/parser"
)

const ethereumRpcURL = "https://ethereum-rpc.publicnode.com"

type EthereumParser struct {
	currentBlock   int
	observers      map[string]bool
	mu             sync.RWMutex
	jsonRPCRequest func(method string, params []interface{}) ([]byte, error)
}

// NewEthereumParser creates a new parser instance
func NewEthereumParser() *EthereumParser {
	return &EthereumParser{
		observers: make(map[string]bool),
		jsonRPCRequest: func(method string, params []interface{}) ([]byte, error) {
			// Default implementation of jsonRPCRequest
			payload := map[string]interface{}{
				"jsonrpc": "2.0",
				"method":  method,
				"params":  params,
				"id":      1,
			}

			jsonPayload, _ := json.Marshal(payload)
			resp, err := http.Post(ethereumRpcURL, "application/json", bytes.NewBuffer(jsonPayload))
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			return io.ReadAll(resp.Body)
		},
	}
}

// GetCurrentBlock retrieves the latest block number
func (p *EthereumParser) GetCurrentBlock() int {
	resp, err := p.jsonRPCRequest("eth_blockNumber", []interface{}{})
	if err != nil {
		log.Printf("Error getting block number: %v", err)
		return p.currentBlock
	}

	var result map[string]string
	json.Unmarshal(resp, &result)

	blockHex := strings.TrimPrefix(result["result"], "0x")
	blockNum, success := new(big.Int).SetString(blockHex, 16)
	if !success {
		log.Printf("Failed to parse block number: %s", blockHex)
		return p.currentBlock
	}

	p.currentBlock = int(blockNum.Int64())
	return p.currentBlock
}

// Subscribe adds an address to observers
func (p *EthereumParser) Subscribe(address string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.observers[strings.ToLower(address)] = true
	return true
}

// GetTransactions retrieves transactions for a given address
func (p *EthereumParser) GetTransactions(address string) []parser.Transaction {
	address = strings.ToLower(address)

	// Ensure the address is subscribed
	if !p.observers[address] {
		return []parser.Transaction{}
	}

	var transactions []parser.Transaction

	// Iterate over the last N blocks (e.g., 10 blocks)
	currentBlock := p.GetCurrentBlock()
	for blockNumber := currentBlock; blockNumber > currentBlock-10 && blockNumber >= 0; blockNumber-- {
		blockHex := "0x" + strconv.FormatInt(int64(blockNumber), 16)

		// Fetch block data with transactions
		resp, err := p.jsonRPCRequest("eth_getBlockByNumber", []interface{}{blockHex, true})

		if err != nil {
			log.Printf("Error fetching block %d: %v", blockNumber, err)
			continue
		}

		// Parse the block response
		var blockData struct {
			Result struct {
				Transactions []parser.Transaction `json:"transactions"`
			} `json:"result"`
		}

		if err := json.Unmarshal(resp, &blockData); err != nil {
			log.Printf("Error parsing transactions for block %d: %v", blockNumber, err)
			continue
		}

		// Filter transactions for the subscribed address
		for _, tx := range blockData.Result.Transactions {
			if strings.EqualFold(tx.From, address) || strings.EqualFold(tx.To, address) {
				transactions = append(transactions, tx)
			}
		}
	}

	return transactions
}

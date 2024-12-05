package ethereum

import (
	"testing"
)

func TestGetCurrentBlock(t *testing.T) {
	parser := NewEthereumParser()

	// Mock jsonRPCRequest to simulate a response
	parser.jsonRPCRequest = func(method string, params []interface{}) ([]byte, error) {
		return []byte(`{"jsonrpc": "2.0", "result": "0x1a"} `), nil // Block 26 in hex
	}

	block := parser.GetCurrentBlock()
	if block != 26 {
		t.Errorf("Expected block number 26, got %d", block)
	}
}

func TestSubscribe(t *testing.T) {
	parser := NewEthereumParser()

	address := "0x123456789abcdef"
	subscribed := parser.Subscribe(address)

	if !subscribed {
		t.Errorf("Expected subscription to return true, got %v", subscribed)
	}

	// Ensure address is added to observers
	if !parser.observers[address] {
		t.Errorf("Address %s was not added to observers", address)
	}
}

func TestGetTransactions(t *testing.T) {
	parser := NewEthereumParser()

	address := "0x123456789abcdef"
	parser.Subscribe(address)

	// Mock jsonRPCRequest to simulate transactions in a block
	parser.jsonRPCRequest = func(method string, params []interface{}) ([]byte, error) {
		return []byte(`{
			"jsonrpc": "2.0",
			"result": {
				"transactions": [
					{"from": "0x123456789abcdef", "to": "0xaabbccddeeff", "value": "100"},
					{"from": "0xaabbccddeeff", "to": "0x123456789abcdef", "value": "200"}
				]
			}
		}`), nil
	}

	transactions := parser.GetTransactions(address)

	if len(transactions) != 2 {
		t.Errorf("Expected 2 transactions, got %d", len(transactions))
	}

	if transactions[0].From != address {
		t.Errorf("Expected transaction 'from' to match %s, got %s", address, transactions[0].From)
	}

	if transactions[1].To != address {
		t.Errorf("Expected transaction 'to' to match %s, got %s", address, transactions[1].To)
	}
}

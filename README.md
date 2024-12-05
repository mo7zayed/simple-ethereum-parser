# Ethereum Transaction Parser

## Overview

This project is a lightweight Ethereum blockchain transaction parser that allows users to:
- Retrieve the current Ethereum block number
- Subscribe to specific Ethereum addresses
- Fetch transactions for subscribed addresses

## Project Structure

```
|-Trustwallet.postman_collection.json
|-main.go
|-README.md
|-go.mod
|-parser
| |-types.go
| |-ethereum
| | |-ethereum_parser_test.go
| | |-ethereum_parser.go
| |-parser.go
|-main_test.go
```

## Installation

1. Clone the repository
```bash
git clone https://github.com/yourusername/trustwallet.git
cd trustwallet
```

2. Install dependencies
```bash
go mod tidy
```

3. Run the application
```bash
go run main.go
```

## API Endpoints

### 1. Get Current Block
- **Endpoint:** `/current_block`
- **Method:** GET
- **Description:** Retrieves the latest Ethereum block number
- **Response:**
  ```json
  {
    "current_block": 12345678
  }
  ```

### 2. Subscribe to Address
- **Endpoint:** `/subscribe`
- **Method:** POST
- **Description:** Add an Ethereum address to the observer list
- **Request Body:**
  ```json
  {
    "address": "0x46340b20830761efd32832a74d7169b29feb9758"
  }
  ```
- **Response:**
  ```json
  {
    "message": "Address subscribed successfully"
  }
  ```

### 3. Get Transactions
- **Endpoint:** `/transactions`
- **Method:** GET
- **Description:** Fetch transactions for a subscribed address
- **Query Parameter:** `address`
- **Response:**
  ```json
  {
    "transactions": [
      {
        "from": "0x...",
        "to": "0x...",
        "value": "..."
      }
    ]
  }
  ```

## Configuration

- RPC Endpoint: Configured in `ethereum_parser.go`
  - Currently using: `https://ethereum-rpc.publicnode.com`

## Testing

Run tests using:
```bash
go test ./...
```

## Postman Collection

A Postman collection is included (`Trustwallet.postman_collection.json`) for easy API testing and documentation.

## Performance Considerations

- Implement connection pooling for RPC requests
- Add caching mechanisms
- Consider rate limiting and error handling for RPC interactions

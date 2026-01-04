# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Gopholler is a Go wrapper for the Telstra Messaging API v3.x. It provides a clean, idiomatic Go interface for sending SMS/MMS messages via Telstra's network in Australia.

## Build Commands

```bash
go build ./...          # Build the project
go test ./...           # Run all tests
go test -v ./...        # Run tests with verbose output
go test -run TestName . # Run a single test
go fmt ./...            # Format code
go vet ./...            # Check for common issues
```

## Architecture

The package is organized by API domain:

| File | Purpose |
|------|---------|
| `client.go` | Main `Client` struct, HTTP request helpers, query builders |
| `auth.go` | OAuth2 token management with thread-safe auto-refresh |
| `types.go` | All request/response structs and type constants |
| `errors.go` | Sentinel errors and `APIError` type for `errors.Is()` support |
| `messages.go` | Send, get, update, delete messages |
| `virtual_numbers.go` | Manage virtual numbers and opt-outs |
| `free_trial.go` | Free trial number management |
| `reports.go` | Report generation and retrieval |
| `sender_names.go` | Get approved sender names (paid accounts) |
| `logs.go` | API call logs |
| `health.go` | Health check endpoint |

## Key Patterns

**Client initialization:**
```go
client := gopholler.NewClient(clientID, clientSecret)
```

**All methods accept context.Context as first parameter:**
```go
msg, err := client.SendMessage(ctx, gopholler.SendMessageRequest{...})
```

**Error handling uses sentinel errors:**
```go
if errors.Is(err, gopholler.ErrRateLimit) {
    // Handle rate limiting
}
```

**OAuth2 tokens are automatically cached and refreshed** - no manual token management needed.

## API Base URLs

- Main API: `https://products.api.telstra.com/messaging/v3`
- OAuth: `https://products.api.telstra.com/v2`

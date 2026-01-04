<p align="center">
  <img src="gopholler-logo.svg" alt="Gopholler Logo" width="200">
</p>

# Gopholler

[![Go](https://github.com/Its-donkey/gopholler/actions/workflows/ci.yml/badge.svg)](https://github.com/Its-donkey/gopholler/actions/workflows/ci.yml)
[![Coverage](https://img.shields.io/badge/coverage-88.6%25-brightgreen)](https://github.com/Its-donkey/gopholler)
[![Go Report Card](https://goreportcard.com/badge/github.com/Its-donkey/gopholler)](https://goreportcard.com/report/github.com/Its-donkey/gopholler)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

A Go wrapper for the [Telstra Messaging API v3.x](https://dev.telstra.com/apis/messaging-api).

## Installation

```bash
go get github.com/Its-donkey/gopholler
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/Its-donkey/gopholler"
)

func main() {
    // Create client with your Telstra API credentials
    client := gopholler.NewClient("your-client-id", "your-client-secret")

    ctx := context.Background()

    // Send an SMS
    msg, err := client.SendMessage(ctx, gopholler.SendMessageRequest{
        To:             "0412345678",
        From:           "0401234567", // Your virtual number
        MessageContent: "Hello from Gopholler!",
    })
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Message sent: %v\n", msg.MessageId)

    // Send to multiple recipients
    msg, err = client.SendMessage(ctx, gopholler.SendMessageRequest{
        To:             []string{"0412345678", "0487654321"},
        From:           "0401234567",
        MessageContent: "Bulk message",
    })

    // Get message status
    status, err := client.GetMessage(ctx, "message-id")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Status: %s\n", status.Status)
}
```

## API Coverage

| API | Methods |
|-----|---------|
| **Messages** | `SendMessage`, `GetMessages`, `GetMessage`, `UpdateMessage`, `DeleteMessage`, `UpdateMessageTags` |
| **Virtual Numbers** | `GetVirtualNumbers`, `AssignVirtualNumber`, `GetVirtualNumber`, `UpdateVirtualNumber`, `DeleteVirtualNumber`, `GetVirtualNumberOptouts` |
| **Free Trial Numbers** | `GetFreeTrialNumbers`, `AddFreeTrialNumbers` |
| **Reports** | `GetReports`, `GetReport`, `CreateMessagesReport` |
| **Sender Names** | `GetSenderNames` |
| **Logs** | `GetLogs` |
| **Health** | `HealthCheck` |

## Error Handling

Errors can be checked using `errors.Is()`:

```go
import "errors"

msg, err := client.SendMessage(ctx, req)
if err != nil {
    switch {
    case errors.Is(err, gopholler.ErrAuthentication):
        // Invalid credentials or expired token
    case errors.Is(err, gopholler.ErrRateLimit):
        // Too many requests
    case errors.Is(err, gopholler.ErrBadRequest):
        // Invalid request parameters
    case errors.Is(err, gopholler.ErrNotFound):
        // Resource not found
    case errors.Is(err, gopholler.ErrPaymentRequired):
        // Account/billing issue
    case errors.Is(err, gopholler.ErrForbidden):
        // Insufficient permissions
    case errors.Is(err, gopholler.ErrServer):
        // Server error
    }
}
```

Access detailed error information:

```go
var apiErrs *gopholler.APIErrors
if errors.As(err, &apiErrs) {
    for _, e := range apiErrs.Errors {
        fmt.Printf("Code: %s, Issue: %s, Suggested: %s\n",
            e.Code, e.Issue, e.SuggestedAction)
    }
}
```

## Client Options

```go
// Custom HTTP client
client := gopholler.NewClient(id, secret,
    gopholler.WithHTTPClient(&http.Client{Timeout: 60 * time.Second}),
)

// Custom timeout
client := gopholler.NewClient(id, secret,
    gopholler.WithTimeout(60 * time.Second),
)

// Custom scopes (if you don't need all permissions)
client := gopholler.NewClient(id, secret,
    gopholler.WithScopes([]string{"messages:read", "messages:write"}),
)
```

## Authentication

OAuth2 tokens are automatically managed. The client will:
- Fetch a token on the first API call
- Cache the token for reuse
- Automatically refresh before expiry

To manually invalidate the cached token:

```go
client.InvalidateToken()
```

## License

MIT

package gopholler

import "time"

// MessageStatus represents the status of a message.
type MessageStatus string

const (
	MessageStatusQueued        MessageStatus = "queued"
	MessageStatusSent          MessageStatus = "sent"
	MessageStatusDelivered     MessageStatus = "delivered"
	MessageStatusExpired       MessageStatus = "expired"
	MessageStatusUndeliverable MessageStatus = "undeliverable"
)

// MessageDirection represents the direction of a message.
type MessageDirection string

const (
	MessageDirectionOutgoing MessageDirection = "outgoing"
	MessageDirectionIncoming MessageDirection = "incoming"
)

// Multimedia represents multimedia content for sending MMS.
type Multimedia struct {
	Type     string `json:"type"`     // e.g., "image/jpeg", "audio/mp3"
	FileName string `json:"fileName"`
	Payload  string `json:"payload"`  // Base64 encoded content
}

// MultimediaGet represents multimedia content in a received message.
type MultimediaGet struct {
	Type     string `json:"type"`
	FileName string `json:"fileName"`
}

// SendMessageRequest represents the request body for sending a message.
type SendMessageRequest struct {
	To                   interface{} `json:"to"`                             // string or []string
	From                 string      `json:"from"`
	MessageContent       string      `json:"messageContent,omitempty"`
	Multimedia           []Multimedia `json:"multimedia,omitempty"`
	RetryTimeout         int         `json:"retryTimeout,omitempty"`         // 10-1440 minutes
	ScheduleSend         *time.Time  `json:"scheduleSend,omitempty"`
	DeliveryNotification bool        `json:"deliveryNotification,omitempty"`
	StatusCallbackUrl    string      `json:"statusCallbackUrl,omitempty"`
	Tags                 []string    `json:"tags,omitempty"`
}

// MessageSent represents the response after sending a message.
type MessageSent struct {
	MessageId            interface{}     `json:"messageId"` // string or []string for bulk sends
	Status               MessageStatus   `json:"status"`
	To                   interface{}     `json:"to"`        // string or []string
	From                 string          `json:"from"`
	MessageContent       string          `json:"messageContent,omitempty"`
	Multimedia           []MultimediaGet `json:"multimedia,omitempty"`
	RetryTimeout         int             `json:"retryTimeout,omitempty"`
	ScheduleSend         *time.Time      `json:"scheduleSend,omitempty"`
	DeliveryNotification bool            `json:"deliveryNotification,omitempty"`
	StatusCallbackUrl    string          `json:"statusCallbackUrl,omitempty"`
	Tags                 []string        `json:"tags,omitempty"`
}

// MessageGet represents a message retrieved from the API.
type MessageGet struct {
	MessageId            string           `json:"messageId"`
	Status               MessageStatus    `json:"status"`
	CreateTimestamp      *time.Time       `json:"createTimestamp,omitempty"`
	SentTimestamp        *time.Time       `json:"sentTimestamp,omitempty"`
	ReceivedTimestamp    *time.Time       `json:"receivedTimestamp,omitempty"`
	To                   string           `json:"to"`
	From                 string           `json:"from"`
	MessageContent       string           `json:"messageContent,omitempty"`
	Multimedia           []MultimediaGet  `json:"multimedia,omitempty"`
	Direction            MessageDirection `json:"direction"`
	RetryTimeout         int              `json:"retryTimeout,omitempty"`
	ScheduleSend         *time.Time       `json:"scheduleSend,omitempty"`
	DeliveryNotification bool             `json:"deliveryNotification,omitempty"`
	StatusCallbackUrl    string           `json:"statusCallbackUrl,omitempty"`
	QueuePriority        int              `json:"queuePriority,omitempty"` // 1-3
	Tags                 []string         `json:"tags,omitempty"`
}

// UpdateMessageRequest represents the request body for updating a scheduled message.
type UpdateMessageRequest struct {
	To                   string       `json:"to"`
	From                 string       `json:"from"`
	MessageContent       string       `json:"messageContent,omitempty"`
	Multimedia           []Multimedia `json:"multimedia,omitempty"`
	RetryTimeout         int          `json:"retryTimeout,omitempty"`
	ScheduleSend         *time.Time   `json:"scheduleSend,omitempty"`
	DeliveryNotification bool         `json:"deliveryNotification,omitempty"`
	StatusCallbackUrl    string       `json:"statusCallbackUrl,omitempty"`
	Tags                 []string     `json:"tags,omitempty"`
}

// MessageUpdate represents the response after updating a message.
type MessageUpdate struct {
	MessageId            string          `json:"messageId"`
	Status               MessageStatus   `json:"status"`
	To                   string          `json:"to"`
	From                 string          `json:"from"`
	MessageContent       string          `json:"messageContent,omitempty"`
	Multimedia           []MultimediaGet `json:"multimedia,omitempty"`
	RetryTimeout         int             `json:"retryTimeout,omitempty"`
	ScheduleSend         *time.Time      `json:"scheduleSend,omitempty"`
	DeliveryNotification bool            `json:"deliveryNotification,omitempty"`
	StatusCallbackUrl    string          `json:"statusCallbackUrl,omitempty"`
	QueuePriority        int             `json:"queuePriority,omitempty"`
	Tags                 []string        `json:"tags,omitempty"`
}

// GetMessagesOptions represents options for listing messages.
type GetMessagesOptions struct {
	Limit     int              `url:"limit,omitempty"`     // 1-50, default 10
	Offset    int              `url:"offset,omitempty"`
	Direction MessageDirection `url:"direction,omitempty"`
	Status    MessageStatus    `url:"status,omitempty"`
	Filter    string           `url:"filter,omitempty"`    // Filter by tags
	Reverse   bool             `url:"reverse,omitempty"`
	StartTime *time.Time       `url:"startTime,omitempty"`
	EndTime   *time.Time       `url:"endTime,omitempty"`
}

// MessagesResponse represents the response for listing messages.
type MessagesResponse struct {
	Messages []MessageGet `json:"messages"`
	Paging   *Paging      `json:"paging,omitempty"`
}

// Paging represents pagination information.
type Paging struct {
	NextPage     string `json:"nextPage,omitempty"`
	PreviousPage string `json:"previousPage,omitempty"`
	LastPage     string `json:"lastPage,omitempty"`
	TotalCount   int    `json:"totalCount,omitempty"`
}

// VirtualNumber represents a virtual number.
type VirtualNumber struct {
	VirtualNumber    string     `json:"virtualNumber"`
	ReplyCallbackUrl string     `json:"replyCallbackUrl,omitempty"`
	Tags             []string   `json:"tags,omitempty"`
	LastUse          *time.Time `json:"lastUse,omitempty"`
}

// AssignVirtualNumberRequest represents the request to assign a virtual number.
type AssignVirtualNumberRequest struct {
	ReplyCallbackUrl string   `json:"replyCallbackUrl,omitempty"`
	Tags             []string `json:"tags,omitempty"`
}

// UpdateVirtualNumberRequest represents the request to update a virtual number.
type UpdateVirtualNumberRequest struct {
	ReplyCallbackUrl string   `json:"replyCallbackUrl,omitempty"`
	Tags             []string `json:"tags,omitempty"`
}

// VirtualNumbersResponse represents the response for listing virtual numbers.
type VirtualNumbersResponse struct {
	VirtualNumbers []VirtualNumber `json:"virtualNumbers"`
	Paging         *Paging         `json:"paging,omitempty"`
}

// RecipientOptout represents an opted-out recipient.
type RecipientOptout struct {
	OptoutNumber    string     `json:"optoutNumber"`
	CreateTimestamp *time.Time `json:"createTimestamp,omitempty"`
}

// OptoutsResponse represents the response for listing optouts.
type OptoutsResponse struct {
	RecipientOptouts []RecipientOptout `json:"recipientOptouts"`
	Paging           *Paging           `json:"paging,omitempty"`
}

// FreeTrialNumbers represents free trial numbers.
type FreeTrialNumbers struct {
	FreeTrialNumbers []string `json:"freeTrialNumbers"`
}

// AddFreeTrialNumbersRequest represents the request to add free trial numbers.
type AddFreeTrialNumbersRequest struct {
	FreeTrialNumbers []string `json:"freeTrialNumbers"`
}

// ReportStatus represents the status of a report.
type ReportStatus string

const (
	ReportStatusQueued    ReportStatus = "queued"
	ReportStatusCompleted ReportStatus = "completed"
	ReportStatusFailed    ReportStatus = "failed"
)

// Report represents a report.
type Report struct {
	ReportId     string       `json:"reportId"`
	ReportStatus ReportStatus `json:"reportStatus"`
	ReportType   string       `json:"reportType,omitempty"`
	ReportExpiry string       `json:"reportExpiry,omitempty"` // date format
	ReportUrl    string       `json:"reportUrl,omitempty"`
}

// ReportsResponse represents the response for listing reports.
type ReportsResponse struct {
	Reports []Report `json:"reports"`
	Paging  *Paging  `json:"paging,omitempty"`
}

// CreateMessagesReportRequest represents the request to create a messages report.
type CreateMessagesReportRequest struct {
	StartDate         string `json:"startDate"`                   // YYYY-MM-DD
	EndDate           string `json:"endDate"`                     // YYYY-MM-DD
	ReportCallbackUrl string `json:"reportCallbackUrl,omitempty"`
	Filter            string `json:"filter,omitempty"`
}

// ReportCreated represents the response after creating a report.
type ReportCreated struct {
	ReportId          string       `json:"reportId"`
	ReportCallbackUrl string       `json:"reportCallbackUrl,omitempty"`
	ReportStatus      ReportStatus `json:"reportStatus"`
}

// ListOptions represents common options for list endpoints.
type ListOptions struct {
	Limit  int    `url:"limit,omitempty"`  // 1-50
	Offset int    `url:"offset,omitempty"`
	Filter string `url:"filter,omitempty"`
}

// HealthCheckResponse represents the response from the health check endpoint.
type HealthCheckResponse struct {
	Status string `json:"status"`
}

// OAuth2Token represents an OAuth2 token response.
type OAuth2Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   string `json:"expires_in"`
}

package messaging

import (
	"time"
)

type ServiceType string

const (
	AuthService   ServiceType = "auth"
	UserService   ServiceType = "user"
	TripService   ServiceType = "trip"
	SocialService ServiceType = "social"
)

type MessageType string

var AuthTypes = struct {
	CreatedUser MessageType
}{

	CreatedUser: "created_user",
}
var UserTypes = struct {
	UpdatedUser MessageType
}{

	UpdatedUser: "updated_user",
}
var SocialTypes = struct {
	FollowUser MessageType
	BlockUser  MessageType
}{

	FollowUser: "follow_user",
	BlockUser:  "block_user",
}

type Message struct {
	ID          string        `json:"id"`           // Unique message ID
	Type        MessageType   `json:"type"`         // Message type (e.g., "user_created")
	Data        interface{}   `json:"data"`         // Actual message payload
	Created     time.Time     `json:"created"`      // Message creation time
	FromService ServiceType   `json:"from_service"` // Source service
	ToServices  []ServiceType `json:"to_services"`  // ToService ervice (empty for broadcast)
	RetryCount  int           `json:"retry_count"`  // Number of retry attempts
	Priority    int           `json:"priority"`     // Message priority (0-9)
	Headers     Headers       `json:"headers"`      // Custom message headers
	Critical    bool          `json:"critical"`
}

type Headers map[string]interface{}

type MessageHandler func(Message) error

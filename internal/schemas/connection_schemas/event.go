package connectionschemas

import "time"

const (
	READY_TO_SYNC = "ready_to_sync"

	ACK_AUTHORIZED   = "ack_auth"
	ACK_SYNC = "ack_sync"
	ACK_MESSAGE_DELIVERY = "ack_message_delivery"
)

type Event struct {
	Type    string `json:"type"`
	Payload EventPayload `json:"payload"`

	EmittedAt time.Time `json:"emitted_at"`
}

type EventPayload interface {}
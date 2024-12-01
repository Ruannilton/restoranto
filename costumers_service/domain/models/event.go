package models

import (
	"encoding/json"
	"time"
)

type OutboxMessage struct {
	Id        int
	Name      string
	Message   json.RawMessage
	CreatedAt time.Time
}

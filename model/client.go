package model

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ClientID         uuid.UUID  `json:"client_id"`
	Name             string     `json:"name"`
	Pronouns         string     `json:"pronouns"`
	Over18           bool       `json:"over_18"`
	PreferredContact string     `json:"preferred_contact"`
	PhoneNumber      int64      `json:"phone_number"`
	CreatedAt        *time.Time `json:"created_at"`
}

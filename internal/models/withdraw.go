package models

import "time"

type WithdrawInfo struct {
	UserName    *string    `json:"user,omitempty"`
	OrderID     string     `json:"order"`
	ProcessedAt *time.Time `json:"processed_at,omitempty"`
	Amount      float64    `json:"sum"`
}

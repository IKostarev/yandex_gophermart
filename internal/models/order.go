package models

import "time"

type OrderStatus string

type OrderInfo struct {
	UserName  *string     `json:"user,omitempty"`
	OrderID   string      `json:"number"`
	Order     *string     `json:"order,omitempty"`
	CreatedAt *time.Time  `json:"uploaded_at,omitempty"`
	Status    OrderStatus `json:"status"`
	Accrual   float64     `json:"accrual"`
}

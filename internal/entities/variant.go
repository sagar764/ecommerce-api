package entities

import "time"

type Variant struct {
	ID            string    `json:"id"`
	Name          string    `json:"name,omitempty"`
	MRP           float64   `json:"mrp"`
	DiscountPrice *float64  `json:"discount_price,omitempty"`
	Size          *string   `json:"size,omitempty"`
	Color         *string   `json:"color,omitempty"`
	Quantity      int       `json:"quantity"`
	IsActive      bool      `json:"is_active,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	CreatedBy     string    `json:"created_by,omitempty"`
	UpdatedBy     string    `json:"updated_by,omitempty"`
}

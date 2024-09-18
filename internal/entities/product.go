package entities

import "time"

type Product struct {
	ID            string   `json:"id"`
	Name          string   `json:"name" binding:"required"`
	Description   string   `json:"description,omitempty"`
	ImageURL      string   `json:"image_url,omitempty"`
	ChildVariants []string `json:"child_variants,omitempty"`
}

type ProductDetails struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	ImageURL    string      `json:"image_url"`
	IsActive    bool        `json:"is_active"`
	CreatedAt   CustomTime  `json:"created_at"`
	UpdatedAt   CustomTime  `json:"updated_at"`
	CreatedBy   string      `json:"created_by"`
	UpdatedBy   string      `json:"updated_by"`
	Categories  []*Category `json:"categories,omitempty"`
	Variants    []*Variant  `json:"variants,omitempty"`
}

// Custom time format for PostgreSQL
const timeFormat = "2006-01-02T15:04:05.999999"

// Custom Time type
type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	str := string(b)
	if str == "null" {
		ct.Time = time.Time{}
		return nil
	}

	// Remove quotes around the time string
	str = str[1 : len(str)-1]

	parsedTime, err := time.Parse(timeFormat, str)
	if err != nil {
		return err
	}
	ct.Time = parsedTime
	return nil
}

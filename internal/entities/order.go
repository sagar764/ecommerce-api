package entities

type Order struct {
	ID     string      `json:"id"`
	Status string      `json:"status"`
	Items  []OrderItem `json:"items"`
	Total  float64     `json:"total"`
}

type OrderItem struct {
	VariantID   string  `json:"variant_id"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	VariantName string  `json:"variant_name,omitempty"`
	ProductName string  `json:"product_name,omitempty"`
}

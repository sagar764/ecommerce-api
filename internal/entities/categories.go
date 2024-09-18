package entities

// Category represents a category in the system
type Category struct {
	ID              string   `json:"id"`
	Name            string   `json:"name" binding:"required"`
	ChildCategories []string `json:"child_categories,omitempty" db:"child_categories"`
	Products        []string `json:"products,omitempty" db:"category_products"`
}

type CategoryDetails struct {
	ID              string      `json:"id"`
	Name            string      `json:"name"`
	IsActive        bool        `json:"is_active"`
	CreatedAt       CustomTime  `json:"created_at"`
	UpdatedAt       CustomTime  `json:"updated_at"`
	CreatedBy       string      `json:"created_by"`
	UpdatedBy       string      `json:"updated_by"`
	ChildCategories []*Category `json:"child_categories"`
	Products        []*Product  `json:"products"`
}

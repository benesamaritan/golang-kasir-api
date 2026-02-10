package models

type Products struct {
	ID                  int    `json:"id"`
	Name                string `json:"name"`
	Price               int    `json:"price"`
	Stock               int    `json:"stock"`
	Category            int    `json:"category_id,omitempty"`
	CategoryName        string `json:"category_name,omitempty"`
	CategoryDescription string `json:"category_description,omitempty"`
}

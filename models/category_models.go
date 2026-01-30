package models

// Bagian Kategori
type Category struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ProductCount int    `json:"productCount,omitempty"`
}

package models

type Categories struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Total       int    `json:"total,omitempty"`
}

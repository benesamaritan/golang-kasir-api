package models

// Bagian Kategori
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"nama"`
	Description string `json:"deskripsi"`
}

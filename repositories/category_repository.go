package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// Create
func (repo *CategoryRepository) Create(category *models.Categories) error {
	query := "INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id"
	err := repo.db.QueryRow(query, category.Name, category.Description).Scan(&category.ID)
	return err
}

// Read - GetAll - ambil Kategori
func (repo *CategoryRepository) GetAll() ([]models.Categories, error) {
	query := "SELECT id, name, description FROM categories"
	rows, err := repo.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]models.Categories, 0)
	for rows.Next() {
		var c models.Categories
		err := rows.Scan(&c.ID, &c.Name, &c.Description)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}

// Read - GetByID - ambil Kategori by ID
func (repo *CategoryRepository) GetByID(id int) (*models.Categories, error) {
	query := `
		SELECT
			id,
			name,
			description,
		    (
		    	SELECT COUNT(*)
		    	FROM products
		    	WHERE products.category_id = categories.id
		    )
		FROM categories
		WHERE categories.id = $1
	`

	var p models.Categories
	err := repo.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Description,
		&p.Total,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("Kategori tidak ditemukan")
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

// Update
func (repo *CategoryRepository) Update(category *models.Categories) error {
	query := "UPDATE categories SET name = $1, description = $2, WHERE id = $3"
	result, err := repo.db.Exec(query, category.Name, category.Description, category.ID)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return errors.New("produk tidak ditemukan")
	}
	return nil
}

// Delete
func (repo *CategoryRepository) Delete(id int) error {
	query := "DELETE FROM categories WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Kategori tidak ditemukan")
	}
	return err
}

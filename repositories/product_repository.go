package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"kasir-api/models"
	"strings"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (repo *ProductRepository) GetAll(name string, active *bool) ([]models.Product, error) {
	query := "SELECT id, name, price, stock, category_id, active FROM products"
	whereClauses := []string{}
	var args []interface{}
	argCounter := 1

	if name != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("name ILIKE $%d", argCounter))
		args = append(args, "%"+name+"%")
		argCounter++
	}
	if active != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("active = $%d", argCounter))
		args = append(args, *active)
		argCounter++
	}

	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	rows, err := repo.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]models.Product, 0)
	for rows.Next() {
		var p models.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.Active)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func (repo *ProductRepository) Create(product *models.Product) error {
	query := "INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id"
	err := repo.db.QueryRow(query, product.Name, product.Price, product.Stock, product.CategoryID).Scan(&product.ID)
	return err
}

// GetByID - ambil produk by ID
func (repo *ProductRepository) GetByID(id int) (*models.Product, error) {
	query := `
		SELECT
			products.id AS product_id,
			products.name,
			products.price,
			products.stock,
			categories.id,
			categories.name,
			categories.description 
		FROM products
		LEFT JOIN categories
		ON products.category_id = categories.id
		WHERE products.id = $1
	`
	var p models.Product
	err := repo.db.QueryRow(query, id).Scan(
		&p.ID,
		&p.Name,
		&p.Price,
		&p.Stock,
		&p.CategoryID,
		&p.CategoryName,
		&p.CategoryDescription,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("Produk tidak ditemukan")
	}

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (repo *ProductRepository) Update(product *models.Product) error {
	query := "UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4 WHERE id = $5"
	result, err := repo.db.Exec(query, product.Name, product.Price, product.Stock, product.CategoryID, product.ID)

	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Produk tidak ditemukan")
	}

	return nil
}

func (repo *ProductRepository) Delete(id int) error {
	query := "DELETE FROM products WHERE id = $1"
	result, err := repo.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Produk tidak ditemukan")
	}

	return err
}

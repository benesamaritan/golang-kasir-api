package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"kasir-api/models"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// inisialisasi subtotal -> jumlah total transaksi keseluruhan
	totalAmount := 0
	// inisialisasi modeling transactionDetails -> nanti kita insert ke db
	details := make([]models.TransactionDetail, 0)
	// loop setiap item
	for _, item := range items {
		var productName string
		var price, stock int
		// get product dapet pricing
		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id=$1 FOR UPDATE", item.ProductID).Scan(&productName, &price, &stock)

		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}

		if err != nil {
			return nil, err
		}

		// hitung current total = quantity * pricing
		// ditambahin ke dalam subtotal
		subtotal := item.Quantity * price
		totalAmount += subtotal

		// kurangi jumlah stok
		res, err := tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)

		if stock < item.Quantity {
			return nil, fmt.Errorf(
				"stock product %s is not enough (available %d)",
				productName,
				stock,
			)
		}

		if err != nil {
			return nil, err
		}

		// item nya dimasukkin ke transactionDetails
		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
		rows, _ := res.RowsAffected()
		if rows == 0 {
			return nil, fmt.Errorf("failed to update stock product id %d", item.ProductID)
		}
	}

	// insert transaction
	var transactionID int
	var createdAt time.Time
	err = tx.QueryRow("INSERT INTO transactions (total_amount, created_at) VALUES ($1, NOW()) RETURNING ID, created_at", totalAmount).Scan(&transactionID, &createdAt)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

// insertTransactionDetails
func (repo *TransactionRepository) insertTransactionDetails(ctx context.Context, tx *sql.Tx, transactionID int, details []models.TransactionDetail) error {
	stmt, err := tx.PrepareContext(ctx, "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return fmt.Errorf("failed to prepare transaction detail statement: %w", err)
	}
	defer stmt.Close()

	for _, detail := range details {
		_, err = stmt.ExecContext(ctx, transactionID, detail.ProductID, detail.Quantity, detail.Subtotal)
		if err != nil {
			return fmt.Errorf("failed to insert transaction detail: %w", err)
		}
	}
	return nil
}

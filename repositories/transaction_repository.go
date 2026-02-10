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
	ctx := context.Background()
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
			if err != nil {
				// Consider logging this error
			}
		}
	}()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		detail, err := repo.processItem(ctx, tx, item)
		if err != nil {
			return nil, err
		}
		totalAmount += detail.Subtotal
		details = append(details, *detail)
	}

	transactionID, createdAt, err := repo.insertTransaction(ctx, tx, totalAmount)
	if err != nil {
		return nil, err
	}

	err = repo.insertTransactionDetails(ctx, tx, transactionID, details)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt:   createdAt,
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) processItem(
	ctx context.Context,
	tx *sql.Tx,
	item models.CheckoutItem,
) (*models.TransactionDetail, error) {
	var productName string
	var price, stock int

	err := tx.QueryRowContext(
		ctx,
		"SELECT name, price, stock FROM products WHERE id=$1 FOR UPDATE",
		item.ProductID,
	).Scan(
		&productName,
		&price,
		&stock,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("product id %d not found", item.ProductID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product info: %w", err)
	}

	if stock < item.Quantity {
		return nil, fmt.Errorf("stock product %s is not enough (available %d)", productName, stock)
	}

	subtotal := item.Quantity * price

	res, err := tx.ExecContext(ctx, "UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed to update stock: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("failed to update stock for product id %d", item.ProductID)
	}

	return &models.TransactionDetail{
		ProductID:   item.ProductID,
		ProductName: productName,
		Quantity:    item.Quantity,
		Subtotal:    subtotal,
	}, nil
}

func (repo *TransactionRepository) insertTransaction(
	ctx context.Context,
	tx *sql.Tx,
	totalAmount int,
) (
	int,
	time.Time,
	error,
) {
	var transactionID int
	var createdAt time.Time

	err := tx.QueryRowContext(
		ctx,
		"INSERT INTO transactions (total_amount, created_at) VALUES ($1, NOW()) RETURNING ID, created_at",
		totalAmount,
	).Scan(
		&transactionID,
		&createdAt,
	)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to insert transaction: %w", err)
	}

	return transactionID, createdAt, nil
}

// insertTransactionDetails
func (repo *TransactionRepository) insertTransactionDetails(
	ctx context.Context,
	tx *sql.Tx,
	transactionID int,
	details []models.TransactionDetail,
) error {
	stmt, err := tx.PrepareContext(
		ctx,
		"INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
	)
	if err != nil {
		return fmt.Errorf("failed to prepare transaction detail statement: %w", err)
	}
	defer stmt.Close()

	for _, detail := range details {
		_, err = stmt.ExecContext(
			ctx,
			transactionID,
			detail.ProductID,
			detail.Quantity,
			detail.Subtotal,
		)
		if err != nil {
			return fmt.Errorf("failed to insert transaction detail: %w", err)
		}
	}
	return nil
}

package repositories

import (
	"database/sql"
	"kasir-api/models"
)

type ReportRepository struct {
	db *sql.DB
}

func NewReportRepository(db *sql.DB) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) GetSummary(start, end string) (int, int, error) {
	var revenue, transaksi int
	err := r.db.QueryRow(`
		SELECT COALESCE(SUM(total_amount),0), COUNT(*)
		FROM transactions
		WHERE (created_at AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta')::date BETWEEN $1 AND $2
	`, start, end).Scan(&revenue, &transaksi)

	return revenue, transaksi, err
}

func (r *ReportRepository) GetBestSeller(start, end string) (*models.BestSeller, error) {
	var bestSeller models.BestSeller
	err := r.db.QueryRow(`
		SELECT p.name, SUM(td.quantity) AS total_qty
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE (t.created_at AT TIME ZONE 'UTC' AT TIME ZONE 'Asia/Jakarta')::date BETWEEN $1 AND $2
		GROUP BY p.name
		ORDER BY total_qty DESC
		LIMIT 1
	`, start, end).Scan(&bestSeller.Name, &bestSeller.Qty)

	if err == sql.ErrNoRows {
		return &models.BestSeller{}, nil
	}
	if err != nil {
		return nil, err
	}
	return &bestSeller, nil
}

package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/jacky-htg/billings/entity"
	_ "github.com/lib/pq"
)

type PaymentRepository struct {
	Db     *sql.DB
	Tx     *sql.Tx
	Log    *log.Logger
	Entity *entity.Payment
}

func (u *PaymentRepository) CreatePayment(ctx context.Context) error {
	query := `
		INSERT INTO payments (loan_id, payment_amount, payment_date, week)
		VALUES ($1, $2, $3, $4) RETURNING payment_id
	`
	err := u.Tx.QueryRowContext(ctx, query, u.Entity.LoanID, u.Entity.PaymentAmount, u.Entity.PaymentDate, u.Entity.Week).Scan(&u.Entity.PaymentID)
	if err != nil {
		u.Log.Printf("Error creating payment: %v", err)
		return err
	}
	return nil
}

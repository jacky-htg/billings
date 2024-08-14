package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/jacky-htg/billings/entity"
	_ "github.com/lib/pq"
)

type PaymentScheduleRepository struct {
	Db     *sql.DB
	Tx     *sql.Tx
	Log    *log.Logger
	Entity *entity.PaymentSchedule
}

func (u *PaymentScheduleRepository) CreatePaymentSchedule(ctx context.Context) error {
	query := `
		INSERT INTO payment_schedule (loan_id, week, amount, due_date) 
		VALUES ($1, $2, $3, $4) 
		RETURNING schedule_id
	`
	err := u.Tx.QueryRowContext(ctx, query, u.Entity.LoanID, u.Entity.Week, u.Entity.Amount, u.Entity.DueDate).Scan(&u.Entity.ScheduleID)
	if err != nil {
		u.Log.Printf("Error creating payment schedule: %v", err)
		return err
	}

	return nil
}

func (u *PaymentScheduleRepository) GetPaymentSchedule(ctx context.Context, loanId int) ([]entity.PaymentSchedule, error) {
	query := `
		SELECT week, due_date, amount, status
		FROM payment_schedule
		WHERE loan_id = $1
		ORDER BY week ASC
	`
	rows, err := u.Db.QueryContext(ctx, query, loanId)
	if err != nil {
		u.Log.Printf("Error query payment schedule: %v", err)
		return nil, err
	}
	defer rows.Close()

	var schedule []entity.PaymentSchedule
	for rows.Next() {
		var item entity.PaymentSchedule
		err := rows.Scan(&item.Week, &item.DueDate, &item.Amount, &item.Status)
		if err != nil {
			u.Log.Printf("Error scan payment schedule: %v", err)
			return nil, err
		}
		schedule = append(schedule, item)
	}

	if err = rows.Err(); err != nil {
		u.Log.Printf("Error get payment schedule: %v", err)
		return nil, err
	}

	return schedule, nil
}

func (u *PaymentScheduleRepository) HasUnpaidOnPrevInstallement(ctx context.Context) (bool, error) {
	var hasPrevUnpaid bool
	query := `SELECT true FROM payment_schedule WHERE status = 'Unpaid' AND loan_id=$1 AND week = $2`
	err := u.Db.QueryRowContext(ctx, query, u.Entity.LoanID, (u.Entity.Week - 1)).Scan(&hasPrevUnpaid)
	if err != nil && err != sql.ErrNoRows {
		u.Log.Printf("Query Has Prev Payment: %v", err)
		return hasPrevUnpaid, err
	}

	return hasPrevUnpaid, nil
}

func (u *PaymentScheduleRepository) GetStatusAndAmount(ctx context.Context) error {
	query := `SELECT status, amount FROM payment_schedule WHERE loan_id=$1 AND week = $2`
	err := u.Db.QueryRowContext(ctx, query, u.Entity.LoanID, (u.Entity.Week)).Scan(&u.Entity.Status, &u.Entity.Amount)
	if err != nil {
		u.Log.Printf("Query Get Status And Amount: %v", err)
		return err
	}

	return nil
}

func (u *PaymentScheduleRepository) UpdatePaymentScheduleAsPaid(ctx context.Context) error {
	query := `
		UPDATE payment_schedule 
		SET status = 'Paid' 
		WHERE loan_id = $1 AND week = $2
	`
	rs, err := u.Tx.ExecContext(ctx, query, u.Entity.LoanID, u.Entity.Week)
	if err != nil {
		u.Log.Printf("Error update payment schedule status: %v", err)
		return err
	}

	affected, err := rs.RowsAffected()
	if err != nil {
		u.Log.Printf("Error get affected: %v", err)
		return err
	}

	if affected != 1 {
		u.Log.Printf("Error affected not expected: %v", err)
		return err
	}

	return nil
}

func (u *PaymentScheduleRepository) GetDelinquentStatus(ctx context.Context) (bool, error) {
	var isDelinquent bool
	query := `
		SELECT CASE WHEN COUNT(*) >= 2 THEN true ELSE false END is_delinquent
		FROM payment_schedule 
		WHERE loan_id = $1 AND due_date <= NOW() AND status = 'Unpaid'
	`
	err := u.Db.QueryRowContext(ctx, query, u.Entity.LoanID).Scan(&isDelinquent)
	if err != nil && err != sql.ErrNoRows {
		u.Log.Printf("Error get delinquest status: %v", err)
	}
	return isDelinquent, err
}

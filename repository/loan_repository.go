package repository

import (
	"context"
	"database/sql"
	"log"

	"github.com/jacky-htg/billings/entity"
	_ "github.com/lib/pq"
)

type LoanRepository struct {
	Db     *sql.DB
	Tx     *sql.Tx
	Log    *log.Logger
	Entity *entity.Loan
}

func (u *LoanRepository) CreateLoan(ctx context.Context) error {
	query := `
		INSERT INTO loans (customer_id, principal_amount, term_in_weeks, interest_rate, outstanding_amount) 
		VALUES ($1, $2, $3, $4, $2) 
		RETURNING loan_id
	`
	err := u.Tx.QueryRowContext(ctx, query, u.Entity.CustomerID, u.Entity.PrincipalAmount, u.Entity.TermInWeeks, u.Entity.InterestRate).Scan(&u.Entity.LoanID)
	if err != nil {
		u.Log.Printf("Error creating loan: %v", err)
		return err
	}

	return nil
}

func (u *LoanRepository) GetLoanById(ctx context.Context) error {
	query := `
		SELECT customer_id, principal_amount, term_in_weeks, interest_rate, created_at
		FROM loans
		WHERE loan_id = $1
	`
	err := u.Db.QueryRowContext(ctx, query, u.Entity.LoanID).Scan(
		&u.Entity.CustomerID, &u.Entity.PrincipalAmount, &u.Entity.TermInWeeks,
		&u.Entity.InterestRate, &u.Entity.CreatedAt)
	if err != nil {
		if err != sql.ErrNoRows {
			u.Log.Printf("Error get loan: %v", err)
		}
		return err
	}

	return nil
}

func (u *LoanRepository) IsExist(ctx context.Context) error {
	query := `SELECT loan_id FROM loans WHERE loan_id = $1`
	err := u.Db.QueryRowContext(ctx, query, u.Entity.LoanID).Scan(&u.Entity.LoanID)
	if err != nil {
		if err != sql.ErrNoRows {
			u.Log.Printf("Error get loan: %v", err)
		}
		return err
	}

	return nil
}

func (u *LoanRepository) UpdateOutstandingAmount(ctx context.Context) error {
	query := `
		UPDATE loans 
		SET outstanding_amount = (outstanding_amount - (principal_amount/term_in_weeks)) 
		WHERE loan_id = $1
	`
	rs, err := u.Tx.ExecContext(ctx, query, u.Entity.LoanID)
	if err != nil {
		u.Log.Printf("Error update outstanding amount: %v", err)
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

func (u *LoanRepository) GetOutstandingAmount(ctx context.Context) error {
	query := `SELECT outstanding_amount FROM loans WHERE loan_id = $1`
	err := u.Db.QueryRowContext(ctx, query, u.Entity.LoanID).Scan(&u.Entity.OutstandingAmount)
	if err != nil && err != sql.ErrNoRows {
		u.Log.Printf("Error get outstanding amount: %v", err)
	}
	return err
}

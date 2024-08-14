package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/jacky-htg/billings/entity"
)

type CustomerRepository struct {
	Db     *sql.DB
	Log    *log.Logger
	Entity *entity.Customer
}

func (u *CustomerRepository) CreateCustomer(ctx context.Context) error {
	query := `INSERT INTO customers (name, email, phone) VALUES ($1, $2, $3) RETURNING customer_id`
	err := u.Db.QueryRowContext(ctx, query, u.Entity.Name, u.Entity.Email, u.Entity.Phone).Scan(&u.Entity.CustomerID)
	if err != nil {
		u.Log.Printf("Error creating customer: %v", err)
		return err
	}
	return nil
}

func (u *CustomerRepository) GetCustomerById(ctx context.Context) error {
	query := `SELECT name, email, phone FROM customers WHERE customer_id = $1`
	err := u.Db.QueryRowContext(ctx, query, u.Entity.CustomerID).Scan(
		&u.Entity.Name, &u.Entity.Email, &u.Entity.Phone,
	)
	if err != nil {
		if err != sql.ErrNoRows {
			u.Log.Printf("Error get customer: %v", err)
		}
		return err
	}
	return nil
}

func (u *CustomerRepository) IsCustomerExist(ctx context.Context) (bool, error) {
	isExist := false
	query := `SELECT EXISTS(SELECT 1 FROM customers WHERE customer_id=$1)`
	err := u.Db.QueryRowContext(ctx, query, u.Entity.CustomerID).Scan(&isExist)
	if err != nil {
		u.Log.Printf("Customer ID %d does not exist", u.Entity.CustomerID)
		return isExist, errors.New("invalid customer ID")
	}

	return isExist, nil
}

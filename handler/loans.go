package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jacky-htg/loans/dto"
	"github.com/jacky-htg/loans/entity"
	"github.com/jacky-htg/loans/repository"
	"github.com/julienschmidt/httprouter"
)

type LoanHandler struct {
	DB  *sql.DB
	Log *log.Logger
}

func (u *LoanHandler) CreateLoan(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()
	var loanReq dto.CreateLoanRequest

	// Decode the request body into LoanRequest struct
	if err := json.NewDecoder(r.Body).Decode(&loanReq); err != nil {
		u.Log.Printf("Error decode request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the request
	if loanReq.CustomerID <= 0 || loanReq.PrincipalAmount <= 0 || loanReq.TermInWeeks <= 0 || loanReq.InterestRate <= 0 {
		http.Error(w, "Invalid loan data", http.StatusBadRequest)
		return
	}

	customerRepo := repository.CustomerRepository{Db: u.DB, Log: u.Log, Entity: &entity.Customer{CustomerID: loanReq.CustomerID}}
	isExist, err := customerRepo.IsCustomerExist(ctx)
	if err != nil || !isExist {
		http.Error(w, "Invalid customer ID", http.StatusBadRequest)
		return
	}

	tx, err := u.DB.Begin()
	if err != nil {
		u.Log.Printf("Failed to start transaction: %v", err)
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}

	// Create the loan in the database
	loanEntity := loanReq.ToEntity()
	repo := repository.LoanRepository{Tx: tx, Log: u.Log, Entity: &loanEntity}
	err = repo.CreateLoan(ctx)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to create loan", http.StatusInternalServerError)
		return
	}

	// Insert the payment schedule
	weeklyAmount := (repo.Entity.PrincipalAmount * (1 + (repo.Entity.InterestRate / 100))) / float64(repo.Entity.TermInWeeks)
	for week := 1; week <= repo.Entity.TermInWeeks; week++ {
		dueDate := time.Now().AddDate(0, 0, 7*week) // due every week
		paymentScheduleEntity := entity.PaymentSchedule{
			LoanID:  repo.Entity.LoanID,
			Week:    week,
			Amount:  weeklyAmount,
			DueDate: dueDate.Format("2006-01-02"),
		}
		paymentScheduleRepo := repository.PaymentScheduleRepository{Tx: tx, Log: u.Log, Entity: &paymentScheduleEntity}
		err = paymentScheduleRepo.CreatePaymentSchedule(ctx)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to create payment schedule", http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()

	// Respond with the created loan ID
	var response dto.CreateLoanResponse
	response.FromEntity(repo.Entity)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (u *LoanHandler) GetLoanById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	loanId := ps.ByName("loanId")

	// Convert loanId to integer
	id, err := strconv.Atoi(loanId)
	if err != nil {
		u.Log.Printf("Error convert id: %v", err)
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	// Get loan details from the database
	repo := repository.LoanRepository{Db: u.DB, Log: u.Log, Entity: &entity.Loan{LoanID: id}}
	err = repo.GetLoanById(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Loan not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Respond with the loan details
	var response dto.GetLoanResponse
	response.FromEntity(repo.Entity)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (u *LoanHandler) GetOutstanding(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	loanId := ps.ByName("loanId")

	// Convert loanId to integer
	id, err := strconv.Atoi(loanId)
	if err != nil {
		u.Log.Printf("Error convert id: %v", err)
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	// Get outstanding amount from the database
	repo := repository.LoanRepository{Db: u.DB, Log: u.Log, Entity: &entity.Loan{LoanID: id}}
	err = repo.GetOutstandingAmount(ctx)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Loan not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Respond with the outstanding amount
	var response dto.GetOutstandingLoanResponse
	response.FromEntity(repo.Entity)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

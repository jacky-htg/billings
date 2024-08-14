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

type PaymentHandler struct {
	DB  *sql.DB
	Log *log.Logger
}

func (u *PaymentHandler) MakePayment(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	loanId := ps.ByName("loanId")

	// Convert loanId to integer
	id, err := strconv.Atoi(loanId)
	if err != nil {
		u.Log.Printf("Error convert id: %v", err)
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	// Decode the incoming JSON request
	var paymentRequest dto.CreatePaymentRequest
	err = json.NewDecoder(r.Body).Decode(&paymentRequest)
	if err != nil {
		u.Log.Printf("Error decode request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Parse the payment date
	_, err = time.Parse("2006-01-02", paymentRequest.Date)
	if err != nil {
		u.Log.Printf("Error parsing date: %v", err)
		http.Error(w, "Invalid date format. Use YYYY-MM-DD.", http.StatusBadRequest)
		return
	}

	{
		loanRepo := repository.LoanRepository{Db: u.DB, Log: u.Log,
			Entity: &entity.Loan{
				LoanID: id,
			},
		}

		err = loanRepo.IsExist(ctx)
		if err == sql.ErrNoRows {
			http.Error(w, "Loan not found", http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, "Failed to process payment", http.StatusInternalServerError)
			return
		}
	}

	{
		paymentScheduleRepo := repository.PaymentScheduleRepository{Db: u.DB, Log: u.Log,
			Entity: &entity.PaymentSchedule{
				LoanID: id,
				Week:   paymentRequest.InstallmentNumber,
			},
		}

		// check prev installment number has paid
		if paymentRequest.InstallmentNumber > 1 {
			hasPrevUnpaid, err := paymentScheduleRepo.HasUnpaidOnPrevInstallement(ctx)
			if err != nil {
				http.Error(w, "Failed to process payment", http.StatusInternalServerError)
				return
			}

			if hasPrevUnpaid {
				http.Error(w, "Has unpaid on prev installment", http.StatusBadRequest)
				return
			}
		}

		// Check Status and Amount
		err = paymentScheduleRepo.GetStatusAndAmount(ctx)
		if err != nil {
			http.Error(w, "Failed to process payment", http.StatusInternalServerError)
			return
		}

		if paymentScheduleRepo.Entity.Status == "Paid" {
			http.Error(w, "Installment already paid", http.StatusBadRequest)
			return
		}

		if paymentScheduleRepo.Entity.Amount != paymentRequest.Amount {
			http.Error(w, "Amount not match", http.StatusBadRequest)
			return
		}
	}

	tx, err := u.DB.Begin()
	if err != nil {
		u.Log.Printf("Failed to start transaction: %v", err)
		http.Error(w, "Failed to process payment", http.StatusInternalServerError)
		return
	}

	// Insert payment into the database
	paymentEntity := paymentRequest.ToEntity()
	paymentEntity.LoanID = id
	repo := repository.PaymentRepository{Tx: tx, Log: u.Log, Entity: &paymentEntity}
	err = repo.CreatePayment(ctx)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to process payment", http.StatusInternalServerError)
		return
	}

	paymentScheduleRepo := repository.PaymentScheduleRepository{Tx: tx, Log: u.Log,
		Entity: &entity.PaymentSchedule{
			LoanID: repo.Entity.LoanID,
			Week:   repo.Entity.Week,
		},
	}
	err = paymentScheduleRepo.UpdatePaymentScheduleAsPaid(ctx)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to process payment", http.StatusInternalServerError)
		return
	}

	loanRepo := repository.LoanRepository{Tx: tx, Log: u.Log,
		Entity: &entity.Loan{
			LoanID: repo.Entity.LoanID,
		},
	}
	err = loanRepo.UpdateOutstandingAmount(ctx)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to process payment", http.StatusInternalServerError)
		return
	}

	tx.Commit()

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	var response dto.CreatePaymentResponse
	response.FromEntity(repo.Entity)
	json.NewEncoder(w).Encode(response)
}

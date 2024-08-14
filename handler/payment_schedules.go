package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/jacky-htg/billings/dto"
	"github.com/jacky-htg/billings/entity"
	"github.com/jacky-htg/billings/repository"
	"github.com/julienschmidt/httprouter"
)

type PaymentScheduleHandler struct {
	DB  *sql.DB
	Log *log.Logger
}

func (u *PaymentScheduleHandler) GetPaymentSchedule(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	loanId := ps.ByName("loanId")

	// Convert loanId to integer
	id, err := strconv.Atoi(loanId)
	if err != nil {
		u.Log.Printf("Error decode request: %v", err)
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	// Get payment schedule from the database
	repo := repository.PaymentScheduleRepository{Db: u.DB, Log: u.Log}
	schedules, err := repo.GetPaymentSchedule(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Loan not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	var response []dto.GetPaymentScheduleResponse
	for _, v := range schedules {
		var tempResponse dto.GetPaymentScheduleResponse
		tempResponse.FromEntity(&v)
		response = append(response, tempResponse)
	}
	json.NewEncoder(w).Encode(response)
}

func (u *PaymentScheduleHandler) IsDelinquent(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	loanId := ps.ByName("loanId")

	// Convert loanId to integer
	id, err := strconv.Atoi(loanId)
	if err != nil {
		u.Log.Printf("Error convert id: %v", err)
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	// Get delinquent status from the database
	repo := repository.PaymentScheduleRepository{Db: u.DB, Log: u.Log, Entity: &entity.PaymentSchedule{LoanID: id}}
	isDelinquent, err := repo.GetDelinquentStatus(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Loan not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Respond with the delinquent status
	response := map[string]interface{}{
		"loan_id":       loanId,
		"is_delinquent": isDelinquent,
	}
	json.NewEncoder(w).Encode(response)
}

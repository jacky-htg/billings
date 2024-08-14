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

type CustomerHandler struct {
	DB  *sql.DB
	Log *log.Logger
}

func (u *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	var customerRequest dto.CreateCustomerRequest

	// Decode request body to Customer struct
	err := json.NewDecoder(r.Body).Decode(&customerRequest)
	if err != nil {
		u.Log.Printf("Error decode request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	customerEntity := customerRequest.ToEntity()
	repo := repository.CustomerRepository{Db: u.DB, Log: u.Log, Entity: &customerEntity}
	err = repo.CreateCustomer(ctx)
	if err != nil {
		http.Error(w, "Failed to create customer", http.StatusInternalServerError)
		return
	}

	// Respond with the created customer ID
	var response dto.CreateCustomerResponse
	response.FromEntity(repo.Entity)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (u *CustomerHandler) GetCustomerById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()
	customerId := ps.ByName("customerId")

	// Convert customerId to integer
	id, err := strconv.Atoi(customerId)
	if err != nil {
		u.Log.Printf("Error convert id: %v", err)
		http.Error(w, "Invalid loan ID", http.StatusBadRequest)
		return
	}

	// Get customer from the database
	repo := repository.CustomerRepository{Db: u.DB, Log: u.Log, Entity: &entity.Customer{CustomerID: id}}
	err = repo.GetCustomerById(ctx)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Loan not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Respond with the loan details
	var response dto.GetCustomerResponse
	response.FromEntity(repo.Entity)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

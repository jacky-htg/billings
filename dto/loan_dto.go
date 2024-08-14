package dto

import "github.com/jacky-htg/loans/entity"

type CreateLoanRequest struct {
	CustomerID      int     `json:"customer_id"`
	PrincipalAmount float64 `json:"principal_amount"`
	InterestRate    float64 `json:"interest_rate"`
	TermInWeeks     int     `json:"term_in_weeks"`
}

func (u *CreateLoanRequest) ToEntity() entity.Loan {
	return entity.Loan{
		CustomerID:      u.CustomerID,
		PrincipalAmount: u.PrincipalAmount,
		TermInWeeks:     u.TermInWeeks,
		InterestRate:    u.InterestRate,
	}
}

type CreateLoanResponse struct {
	LoanID  int    `json:"loan_id"`
	Message string `json:"message"`
}

func (resp *CreateLoanResponse) FromEntity(e *entity.Loan) {
	resp.LoanID = e.LoanID
	resp.Message = "Loan created successfully"
}

type GetLoanResponse struct {
	LoanID          int     `json:"loan_id"`
	CustomerID      int     `json:"customer_id"`
	PrincipalAmount float64 `json:"principal_amount"`
	TermInWeeks     int     `json:"term_in_weeks"`
	InterestRate    float64 `json:"interest_rate"`
}

func (resp *GetLoanResponse) FromEntity(e *entity.Loan) {
	resp.LoanID = e.LoanID
	resp.CustomerID = e.CustomerID
	resp.PrincipalAmount = e.PrincipalAmount
	resp.TermInWeeks = e.TermInWeeks
	resp.InterestRate = e.InterestRate
}

type GetOutstandingLoanResponse struct {
	LoanID            int     `json:"loan_id"`
	OutstandingAmount float64 `json:"outstanding_amount"`
}

func (resp *GetOutstandingLoanResponse) FromEntity(e *entity.Loan) {
	resp.LoanID = e.LoanID
	resp.OutstandingAmount = e.OutstandingAmount
}

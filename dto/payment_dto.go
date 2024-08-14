package dto

import "github.com/jacky-htg/loans/entity"

type CreatePaymentRequest struct {
	Amount            float64 `json:"amount"`
	Date              string  `json:"date"`
	InstallmentNumber int     `json:"installment_number"`
}

func (u *CreatePaymentRequest) ToEntity() entity.Payment {
	return entity.Payment{
		PaymentAmount: u.Amount,
		PaymentDate:   u.Date,
		Week:          u.InstallmentNumber,
	}
}

type CreatePaymentResponse struct {
	PaymentID int    `json:"payment_id"`
	Message   string `json:"message"`
}

func (resp *CreatePaymentResponse) FromEntity(e *entity.Payment) {
	resp.PaymentID = e.PaymentID
	resp.Message = "Payment successfully recorded"
}

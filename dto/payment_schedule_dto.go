package dto

import "github.com/jacky-htg/billings/entity"

type GetPaymentScheduleResponse struct {
	InstallmentNumber int     `json:"installment_number"`
	DueDate           string  `json:"due_date"`
	Amount            float64 `json:"amount"`
	Status            string  `json:"status"`
}

func (resp *GetPaymentScheduleResponse) FromEntity(e *entity.PaymentSchedule) {
	resp.InstallmentNumber = e.Week
	resp.DueDate = e.DueDate[:10]
	resp.Amount = e.Amount
	resp.Status = e.Status
}

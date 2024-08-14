package entity

type PaymentSchedule struct {
	ScheduleID int
	LoanID     int
	Week       int
	Amount     float64
	DueDate    string
	Status     string
}

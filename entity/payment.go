package entity

type Payment struct {
	PaymentID     int
	LoanID        int
	PaymentAmount float64
	PaymentDate   string
	Week          int
	Status        string
}

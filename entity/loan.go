package entity

type Loan struct {
	LoanID            int
	CustomerID        int
	PrincipalAmount   float64
	InterestRate      float64
	TermInWeeks       int
	OutstandingAmount float64
	CreatedAt         string
}

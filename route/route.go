package route

import (
	"database/sql"
	"log"

	"github.com/jacky-htg/loans/handler"
	"github.com/jacky-htg/loans/middleware"
	"github.com/julienschmidt/httprouter"
)

func InitRoute(db *sql.DB, log *log.Logger) *httprouter.Router {
	router := httprouter.New()

	customer := handler.CustomerHandler{DB: db, Log: log}
	loan := handler.LoanHandler{DB: db, Log: log}
	paymentSchedule := handler.PaymentScheduleHandler{DB: db, Log: log}
	payment := handler.PaymentHandler{DB: db, Log: log}

	mid := middleware.Middleware{DB: db, Log: log}
	middlewares := []middleware.MiddHandler{
		mid.TokenMiddleware, // Apply the token check middleware
	}

	router.POST("/customers", mid.Init(middlewares, customer.CreateCustomer))
	router.GET("/customers/:customerId", mid.Init(middlewares, customer.GetCustomerById))
	router.POST("/loans", mid.Init(middlewares, loan.CreateLoan))
	router.GET("/loans/:loanId", mid.Init(middlewares, loan.GetLoanById))
	router.GET("/loans/:loanId/schedule", mid.Init(middlewares, paymentSchedule.GetPaymentSchedule))
	router.POST("/loans/:loanId/payment", mid.Init(middlewares, payment.MakePayment))
	router.GET("/loans/:loanId/outstanding", mid.Init(middlewares, loan.GetOutstanding))
	router.GET("/loans/:loanId/delinquent", mid.Init(middlewares, paymentSchedule.IsDelinquent))

	return router

}

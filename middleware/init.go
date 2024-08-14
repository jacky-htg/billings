package middleware

import (
	"database/sql"
	"log"

	"github.com/julienschmidt/httprouter"
)

type Middleware struct {
	DB  *sql.DB
	Log *log.Logger
}

type MiddHandler func(httprouter.Handle) httprouter.Handle

func (u *Middleware) Init(mw []MiddHandler, handler httprouter.Handle) httprouter.Handle {
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler)
		}
	}

	return handler
}

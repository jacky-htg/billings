package middleware

import (
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func (u *Middleware) TokenMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		token := r.Header.Get("Authorization")

		// Simple token validation (Bearer token)
		if token == "" || !strings.HasPrefix(token, "Bearer ") || len(token) <= 7 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract the token and validate (add your validation logic here)
		actualToken := token[7:]
		if !isValidToken(actualToken) {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Proceed to the next handler if the token is valid
		next(w, r, ps)
	}
}

func isValidToken(token string) bool {
	// Implement your token validation logic here (e.g., check against a database or a secret)
	return token == "your-secret-token"
}

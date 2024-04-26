package middleware

import (
	"net/http"
	"os"
)

func SecretKeyMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientSecretKey := r.Header.Get("X-Secret-Key")
		if clientSecretKey != os.Getenv("apiSecretKey") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

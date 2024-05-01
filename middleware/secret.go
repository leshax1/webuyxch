package middleware

import (
	"net/http"
	"os"
)

func SecretKeyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientSecretKey := r.Header.Get("X-Secret-Key")
		if clientSecretKey != os.Getenv("apiSecretKey") {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

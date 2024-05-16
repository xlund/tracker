package middleware

import "net/http"

func IsAuthenticated() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// TODO
			next.ServeHTTP(w, r)
		})
	}
}

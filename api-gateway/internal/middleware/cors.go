package middleware

import (
	"net/http"
	"github.com/gorilla/mux"
)

func CORS(allowed map[string]bool) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")
			if origin != "" && allowed[origin] {
				// If you send cookies from the browser, you MUST echo the origin
				// (no "*") and allow credentials.
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				// Add any custom headers your frontend sends (Stripe libs often add one)
				w.Header().Set("Access-Control-Allow-Headers",
					"Content-Type, Authorization, X-Requested-With, Stripe-Version")
			}

			// Preflight short-circuit
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

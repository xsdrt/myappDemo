package middleware

import (
	"net/http"
)

func (m *Middleware) Auth(next http.Handler) http.Handler { //web based auth only, need api also
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !m.App.Session.Exists(r.Context(), "userID") {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	})
}

package middlewares

import (
	"net/http"

	"github.com/roscrl/light/core/views"
)

func RequestTurboStream(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !views.IsTurboStreamRequest(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}

		next.ServeHTTP(w, r)
	})
}

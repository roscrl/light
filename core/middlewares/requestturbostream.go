package middlewares

import (
	"net/http"

	"github.com/roscrl/light/core/views"
)

func RequestTurboStream(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !views.IsTurboStreamRequest(r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}

		next.ServeHTTP(w, r)
	}
}

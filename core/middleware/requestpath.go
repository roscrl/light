package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/roscrl/light/core/util/contextutil"
)

func RequestPath(next http.Handler, ignorePath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if strings.HasPrefix(path, ignorePath) {
			next.ServeHTTP(w, r)

			return
		}

		rctx := context.WithValue(r.Context(), contextutil.RequestPathKey{}, path)
		r = r.WithContext(rctx)

		next.ServeHTTP(w, r)
	})
}

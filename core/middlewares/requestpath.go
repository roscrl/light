package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/roscrl/light/core/utils/contextutil"
)

func RequestPath(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if strings.HasPrefix(path, "/assets") || strings.HasPrefix(path, "/local/browser/refresh") {
			next.ServeHTTP(w, r)

			return
		}

		rctx := context.WithValue(r.Context(), contextutil.RequestPathKey{}, path)
		r = r.WithContext(rctx)

		next.ServeHTTP(w, r)
	}
}

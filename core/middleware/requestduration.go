package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/roscrl/light/core/support/rlog"
	"github.com/roscrl/light/core/support/rlog/key"
)

func RequestDuration(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/assets") {
			next.ServeHTTP(w, r)

			return
		}

		start := time.Now()

		next.ServeHTTP(w, r)

		elapsed := time.Since(start)

		log, rctx := rlog.L(r)
		log.InfoContext(rctx, "request duration", key.Took, elapsed)
	})
}

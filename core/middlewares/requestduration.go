package middlewares

import (
	"net/http"
	"strings"
	"time"

	"github.com/roscrl/light/core/helpers/rlog"
	"github.com/roscrl/light/core/helpers/rlog/key"
)

func RequestDuration(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/assets") || strings.HasPrefix(r.URL.Path, "/local/browser/refresh") {
			next.ServeHTTP(w, r)

			return
		}

		start := time.Now()

		next.ServeHTTP(w, r)

		elapsed := time.Since(start)

		log, rctx := rlog.L(r)
		log.InfoContext(rctx, "request duration", key.Took, elapsed)
	}
}

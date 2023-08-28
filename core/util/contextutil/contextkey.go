package contextutil

import "net/http"

type (
	RequestLoggerKey struct{}
	RequestIDKey     struct{}
	RequestPathKey   struct{}
	RequestIPKey     struct{}
)

// TODO check this
// RequestID returns the request ID from the context otherwise it returns an empty string.
func RequestID(r *http.Request) string {
	return r.Context().Value(RequestIDKey{}).(string)
}

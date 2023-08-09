package contexthelp

import "net/http"

type RequestLogger struct{}

type (
	RequestLoggerKey struct{}
	RequestIDKey     struct{}
	RequestPathKey   struct{}
)

// TODO check this
// RequestID returns the request ID from the context otherwise it returns an empty string.
func RequestID(r *http.Request) string {
	return r.Context().Value(RequestIDKey{}).(string)
}

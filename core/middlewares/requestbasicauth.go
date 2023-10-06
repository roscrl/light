package middlewares

import (
	"crypto/sha256"
	"crypto/subtle"
	"net/http"
)

// https://www.alexedwards.net/blog/basic-authentication-in-go
func RequestBasicAuth(next http.HandlerFunc, correctUsername, correctPassword string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))

			expectedUsernameHash := sha256.Sum256([]byte(correctUsername))
			expectedPasswordHash := sha256.Sum256([]byte(correctPassword))

			usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
			passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)

				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
	})
}

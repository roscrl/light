package app

import (
	"context"
	"net/http"
	"strings"

	"github.com/roscrl/light/core/views"
)

type contextKeyFields struct{}

// https://benhoyt.com/writings/go-routing/
func (app *App) routing(routes []route) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var allow []string

		for _, route := range routes {
			matches := route.regex.FindStringSubmatch(r.URL.Path)
			if len(matches) > 0 {
				if r.Method != route.method {
					allow = append(allow, route.method)

					continue
				}

				ctx := context.WithValue(r.Context(), contextKeyFields{}, matches[1:])
				route.handler(w, r.WithContext(ctx))

				return
			}
		}

		if len(allow) > 0 {
			w.Header().Set("Allow", strings.Join(allow, ", "))
			http.Redirect(w, r, "/", http.StatusSeeOther)

			return
		}

		w.WriteHeader(http.StatusNotFound)
		app.Views.RenderPage(w, views.NotFound, nil)
	})
}

func getField(r *http.Request, index int) string {
	fields := r.Context().Value(contextKeyFields{}).([]string)

	return fields[index]
}

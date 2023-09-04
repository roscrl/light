package app

import (
	"net/http"
)

func (app *App) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := app.DB.Exec("SELECT 1")
		if err != nil {
			app.Views.RenderErrorPage(w, "Hmm, 'SELECT 1' query to DB failed", http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

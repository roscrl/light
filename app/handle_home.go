package app

import (
	"net/http"

	"github.com/roscrl/light/core/helpers/rlog"
	"github.com/roscrl/light/core/helpers/rlog/key"
	"github.com/roscrl/light/core/views"
	"github.com/roscrl/light/core/views/params"
)

func (app *App) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log, rctx := rlog.L(r)

		todos, err := app.Qry.ReadTodos(rctx)
		if err != nil {
			log.ErrorContext(rctx, "failed to query for todos", key.Err, err)
			app.Views.RenderDefaultErrorPage(w)

			return
		}

		app.Views.RenderPage(w, views.Home, map[string]any{
			params.Todos: todos,
		})
	}
}

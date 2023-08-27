package app

import (
	"net/http"

	"github.com/roscrl/light/core/db/sqlc"
	"github.com/roscrl/light/core/domain/todo"
	"github.com/roscrl/light/core/support/rlog"
	"github.com/roscrl/light/core/support/rlog/key"
	"github.com/roscrl/light/core/support/ulid"
	"github.com/roscrl/light/core/views"
	"github.com/roscrl/light/core/views/params"
)

func (app *App) handleTodoCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log, rctx := rlog.L(r)

		if !views.IsTurboStreamRequest(r) {
			http.Redirect(w, r, RouteHome, http.StatusSeeOther)

			return
		}

		if err := r.ParseForm(); err != nil {
			log.InfoContext(rctx, "failed to parse form", key.Err, err)
			app.Views.RenderTurboStream(w, views.TodoFormNewStream, map[string]any{
				params.Error: "Failed to parse given task.",
			})

			return
		}

		task := r.Form.Get("task")
		if task == "" {
			app.Views.RenderTurboStream(w, views.TodoFormNewStream, map[string]any{
				params.Error: "task cannot be empty.",
			})

			return
		}

		todo := sqlc.CreateTodoParams{
			ID:     ulid.NewString(),
			Task:   task,
			Status: "pending",
		}

		_, err := app.Qry.CreateTodo(rctx, todo)
		if err != nil {
			log.ErrorContext(rctx, "failed to create todo", key.Err, err)
			app.Views.RenderTurboStream(w, views.TodoFormNewStream, map[string]any{
				params.Error: "Oops, something went wrong inserting your todo to the database, try again later!",
			})

			return
		}

		app.Views.RenderTurboStream(w, views.TodoFormNewStream, map[string]any{
			params.Todo:      todo,
			params.InputTodo: "",
		})
	}
}

func (app *App) handleTodoEdit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log, rctx := rlog.L(r)

		if !views.IsTurboStreamRequest(r) {
			http.Redirect(w, r, RouteHome, http.StatusSeeOther)

			return
		}

		id := getField(r, 0)

		todo, err := app.Qry.GetTodo(rctx, id)
		if err != nil {
			log.ErrorContext(rctx, "failed to get todo", key.Err, err)
			app.Views.RenderDefaultErrorPage(w)

			return
		}

		log.InfoContext(rctx, "returning todo to edit stream", key.Todo, todo)

		app.Views.RenderTurboStream(w, views.TodoCardEditStream, map[string]any{
			params.Todo: todo,
		})
	}
}

func (app *App) handleTodoUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log, rctx := rlog.L(r)

		if !views.IsTurboStreamRequest(r) {
			http.Redirect(w, r, RouteHome, http.StatusSeeOther)

			return
		}

		if err := r.ParseForm(); err != nil {
			log.InfoContext(rctx, "failed to parse form", key.Err, err)
			app.Views.RenderDefaultErrorPage(w)

			return
		}

		todoID := getField(r, 0)

		existingTodo, err := app.Qry.GetTodo(rctx, todoID)
		if err != nil {
			app.Views.RenderDefaultErrorPage(w)

			return
		}

		status := r.Form.Get("status")
		if status == "" {
			status = string(todo.Pending)
		} else {
			status = string(todo.Done)
		}

		task := r.Form.Get("task")
		if task == "" {
			task = existingTodo.Task
		}

		updatedTodo, err := app.Qry.UpdateTodo(rctx, sqlc.UpdateTodoParams{
			ID:     getField(r, 0),
			Task:   task,
			Status: status,
		})
		if err != nil {
			app.Views.RenderDefaultErrorPage(w)

			return
		}

		app.Views.RenderTurboStream(w, views.TodoCardUpdateStream, map[string]any{
			params.Todo: updatedTodo,
		})
	}
}

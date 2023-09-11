package app

import (
	"net/http"

	"github.com/roscrl/light/core/helpers/rlog"
	"github.com/roscrl/light/core/helpers/rlog/key"
	"github.com/roscrl/light/core/helpers/ulid"
	"github.com/roscrl/light/core/models/todo"
	"github.com/roscrl/light/core/views"
	"github.com/roscrl/light/core/views/params"
	"github.com/roscrl/light/db/sqlc"
)

func (app *App) handleTodosCreate() http.HandlerFunc {
	const formTask = "task"

	return func(w http.ResponseWriter, r *http.Request) {
		log, rctx := rlog.L(r)

		task := r.FormValue(formTask)
		if task == "" {
			app.Views.RenderTurboStream(w, views.TodoFormNewStream, map[string]any{
				params.Error: "task cannot be empty.",
			})

			return
		}

		todo := sqlc.NewTodoParams{
			ID:     ulid.NewString(),
			Task:   task,
			Status: string(todo.Pending),
		}

		_, err := app.Qry.NewTodo(rctx, todo)
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

func (app *App) handleTodosEdit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log, rctx := rlog.L(r)

		id := getField(r, 0)

		todo, err := app.Qry.GetTodoByID(rctx, id)
		if err != nil {
			log.InfoContext(rctx, "failed to get todo", key.ID, id, key.Err, err)
			app.Views.RenderDefaultErrorPage(w)

			return
		}

		app.Views.RenderTurboStream(w, views.TodoCardEditStream, map[string]any{
			params.Todo: todo,
		})
	}
}

func (app *App) handleTodosUpdate() http.HandlerFunc {
	const (
		formStatus = "status"
		formTask   = "task"
	)

	return func(w http.ResponseWriter, r *http.Request) {
		_, rctx := rlog.L(r)

		todoID := getField(r, 0)

		existingTodo, err := app.Qry.GetTodoByID(rctx, todoID)
		if err != nil {
			app.Views.RenderDefaultErrorPage(w)

			return
		}

		status := r.FormValue(formStatus)
		if status == "" {
			status = string(todo.Pending)
		} else {
			status = string(todo.Done)
		}

		task := r.FormValue(formTask)
		if task == "" {
			task = existingTodo.Task
		}

		updatedTodo, err := app.Qry.UpdateTodoByID(rctx, sqlc.UpdateTodoByIDParams{
			ID:     todoID,
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

func (app *App) handleTodosDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//log, rctx := rlog.L(r)
		//
		//todoID := getField(r, 0)

		//if err := jobs.Enqueue(rctx, jobs.TodoDelete, tododelete.Args(todoID), app.Qry); err != nil {
		//	log.ErrorContext(rctx, "failed to enqueue job", key.Err, err)
		//	app.Views.RenderDefaultErrorPage(w)
		//
		//	return
		//}
	}
}

func (app *App) handleTodosSearch() http.HandlerFunc {
	const (
		formQuery = "query"
	)

	return func(w http.ResponseWriter, r *http.Request) {
		log, rctx := rlog.L(r)

		var (
			todos []sqlc.Todo
			err   error
		)

		query := r.FormValue(formQuery)
		if query == "" {
			todos, err = app.Qry.GetAllTodos(rctx)
		} else {
			todos, err = app.Qry.SearchTodosByTask(rctx, "%"+query+"%")
		}

		if err != nil {
			log.ErrorContext(rctx, "failed to search todos", key.Err, err)
			app.Views.RenderTurboStream(w, views.TodoFormSearchStream, map[string]any{
				params.Error: "Oops, something went wrong searching the todos, try again later!",
			})

			return
		}

		app.Views.RenderTurboStream(w, views.TodoListSearchStream, map[string]any{
			params.Todos: todos,
		})
	}
}

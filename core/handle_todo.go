package core

import (
	"net/http"

	"github.com/roscrl/light/core/db/sqlc"
	"github.com/roscrl/light/core/support/rlog"
	"github.com/roscrl/light/core/support/rlog/key"
)

func (s *Server) handleTodoCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log, rctx := rlog.L(r.Context())

		todo := sqlc.CreateTodoParams{
			ID:     "",
			Task:   "test",
			Status: "pending",
		}

		_, err := s.Qry.CreateTodo(rctx, todo)
		if err != nil {
			log.ErrorContext(rctx, "failed to add playlist", key.Err, err)
			s.Views.Stream(w, "playlist/_new.stream.tmpl", map[string]any{
				// "todo_input": playlistLinkOrID,
				"error": "Oops, something went wrong inserting your todo to the database, try again later!",
			})

			return
		}

		s.Views.Render(w, "index.tmpl", map[string]any{
			// "playlists": playlists,
		})
	}
}

package core

import (
	"net/http"

	"github.com/roscrl/light/core/support/rlog"
	"github.com/roscrl/light/core/support/rlog/key"
	"github.com/roscrl/light/core/views"
)

func (s *Server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log, rctx := rlog.L(r.Context())

		log.InfoContext(rctx, "fetching all counts")

		todos, err := s.Qry.GetTodos(rctx)
		if err != nil {
			log.ErrorContext(rctx, "failed to query for todos", key.Err, err)
			s.Views.RenderDefaultError(w)

			return
		}

		s.Views.Render(w, views.Index, map[string]any{
			"todos": todos,
		})
	}
}

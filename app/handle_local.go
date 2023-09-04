package app

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func (app *App) handleLocalBrowserRefresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//nolint:bodyclose
		rc := http.NewResponseController(w)

		if err := rc.SetWriteDeadline(time.Time{}); err != nil {
			log.Fatalf("failed to set write deadline: %v", err)
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-store")

		<-app.Views.LocalBrowserRefreshNotify

		fmt.Fprint(w, "data: refresh\n\n")
		rc.Flush()
	}
}

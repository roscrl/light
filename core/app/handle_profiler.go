package app

import (
	"net/http"
	"net/http/pprof"
)

func (app *App) handleIndex() http.HandlerFunc {
	return pprof.Index
}

func (app *App) handleAllocs() http.HandlerFunc {
	return pprof.Handler("allocs").ServeHTTP
}

func (app *App) handleBlock() http.HandlerFunc {
	return pprof.Handler("block").ServeHTTP
}

func (app *App) handleCmdline() http.HandlerFunc {
	return pprof.Cmdline
}

func (app *App) handleGoroutine() http.HandlerFunc {
	return pprof.Handler("goroutine").ServeHTTP
}

func (app *App) handleHeap() http.HandlerFunc {
	return pprof.Handler("heap").ServeHTTP
}

func (app *App) handleMutex() http.HandlerFunc {
	return pprof.Handler("mutex").ServeHTTP
}

func (app *App) handleProfile() http.HandlerFunc {
	return pprof.Profile
}

func (app *App) handleThreadcreate() http.HandlerFunc {
	return pprof.Handler("threadcreate").ServeHTTP
}

func (app *App) handleSymbol() http.HandlerFunc {
	return pprof.Symbol
}

func (app *App) handleTrace() http.HandlerFunc {
	return pprof.Trace
}

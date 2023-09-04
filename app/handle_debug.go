package app

import (
	"net/http"
	"net/http/pprof"
)

func (app *App) handleDebugIndex() http.HandlerFunc {
	return pprof.Index
}

func (app *App) handleDebugAllocs() http.HandlerFunc {
	return pprof.Handler("allocs").ServeHTTP
}

func (app *App) handleDebugBlock() http.HandlerFunc {
	return pprof.Handler("block").ServeHTTP
}

func (app *App) handleDebugCmdline() http.HandlerFunc {
	return pprof.Cmdline
}

func (app *App) handleDebugGoroutine() http.HandlerFunc {
	return pprof.Handler("goroutine").ServeHTTP
}

func (app *App) handleDebugHeap() http.HandlerFunc {
	return pprof.Handler("heap").ServeHTTP
}

func (app *App) handleDebugMutex() http.HandlerFunc {
	return pprof.Handler("mutex").ServeHTTP
}

func (app *App) handleDebugProfile() http.HandlerFunc {
	return pprof.Profile
}

func (app *App) handleDebugThreadcreate() http.HandlerFunc {
	return pprof.Handler("threadcreate").ServeHTTP
}

func (app *App) handleDebugSymbol() http.HandlerFunc {
	return pprof.Symbol
}

func (app *App) handleDebugTrace() http.HandlerFunc {
	return pprof.Trace
}

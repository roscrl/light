package app

import (
	"net/http"
	"regexp"

	"github.com/roscrl/light/config"
	"github.com/roscrl/light/core/middlewares"
)

const (
	RouteAssetsBase = "/assets"
	RouteAssets     = "/assets/(.*)"

	RouteHome        = "/"
	RouteTodosCreate = "/todos/create"
	RouteTodosEdit   = "/todos/(.*)/edit"
	RouteTodosUpdate = "/todos/(.*)/update"
	RouteTodosSearch = "/todos/search(.*)"

	RouteHealth            = "/health"
	RouteDebugPprof        = "/debug/pprof"
	RouteDebugAllocs       = "/debug/allocs"
	RouteDebugBlock        = "/debug/block"
	RouteDebugCmdline      = "/debug/cmdline"
	RouteDebugGoroutine    = "/debug/goroutine"
	RouteDebugHeap         = "/debug/heap"
	RouteDebugMutex        = "/debug/mutex"
	RouteDebugProfile      = "/debug/profile"
	RouteDebugThreadcreate = "/debug/threadcreate"
	RouteDebugSymbol       = "/debug/symbol"
	RouteDebugTrace        = "/debug/trace"

	RouteLocalBrowserRefresh = "/local/browser/refresh"
)

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

func (app *App) routes() http.Handler {
	newRoute := func(method, pattern string, handler http.HandlerFunc) route {
		return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
	}

	routes := []route{
		newRoute(http.MethodGet, RouteAssets, app.handleAssets()),
		newRoute(http.MethodGet, RouteHome, app.handleHome()),
	}

	{
		turboStreamRoutes := []route{
			newRoute(http.MethodPost, RouteTodosCreate, app.handleTodosCreate()),
			newRoute(http.MethodGet, RouteTodosEdit, app.handleTodosEdit()),
			newRoute(http.MethodPost, RouteTodosUpdate, app.handleTodosUpdate()),
			newRoute(http.MethodGet, RouteTodosSearch, app.handleTodosSearch()),
		}

		for _, route := range turboStreamRoutes {
			handler := middlewares.RequestTurboStream(route.handler)
			route.handler = handler.ServeHTTP

			routes = append(routes, route)
		}
	}

	{
		observabilityRoutes := []route{
			newRoute(http.MethodGet, RouteHealth, app.handleHealth()),

			newRoute(http.MethodGet, RouteDebugPprof, app.handleDebugIndex()),
			newRoute(http.MethodGet, RouteDebugAllocs, app.handleDebugAllocs()),
			newRoute(http.MethodGet, RouteDebugBlock, app.handleDebugBlock()),
			newRoute(http.MethodGet, RouteDebugCmdline, app.handleDebugCmdline()),
			newRoute(http.MethodGet, RouteDebugGoroutine, app.handleDebugGoroutine()),
			newRoute(http.MethodGet, RouteDebugHeap, app.handleDebugHeap()),
			newRoute(http.MethodGet, RouteDebugMutex, app.handleDebugMutex()),
			newRoute(http.MethodGet, RouteDebugProfile, app.handleDebugProfile()),
			newRoute(http.MethodGet, RouteDebugThreadcreate, app.handleDebugThreadcreate()),
			newRoute(http.MethodGet, RouteDebugSymbol, app.handleDebugSymbol()),
			newRoute(http.MethodGet, RouteDebugTrace, app.handleDebugTrace()),
		}

		routes = append(routes, observabilityRoutes...)
	}

	if app.Cfg.Env == config.LOCAL {
		routes = append(routes, newRoute(http.MethodGet, RouteLocalBrowserRefresh, app.handleLocalBrowserRefresh()))
	}

	middlewares := []func(http.Handler) http.HandlerFunc{
		middlewares.RequestLogger,
		middlewares.RequestPath,
		middlewares.RequestID,
		middlewares.RequestIP,
		middlewares.RequestDuration,
		middlewares.RequestRecoverer,
	}

	routerEntry := app.routing(routes)
	wrappedEntry := http.Handler(routerEntry)

	for i := len(middlewares) - 1; i >= 0; i-- {
		mw := middlewares[i]
		wrappedEntry = mw(wrappedEntry)
	}

	return wrappedEntry
}

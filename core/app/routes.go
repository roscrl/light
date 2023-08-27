package app

import (
	"net/http"
	"regexp"

	"github.com/roscrl/light/core/middleware"
)

const (
	RouteAssetBase = "/assets"
	RouteAsset     = "/assets/(.*)"

	RouteHome       = "/"
	RouteTodoCreate = "/todos/create"
	RouteTodoEdit   = "/todos/(.*)/edit"
	RouteTodoUpdate = "/todos/(.*)/update"

	RouteHealth              = "/health"
	RouteProfileBaseRoute    = "/debug/pprof"
	RouteProfileAllocs       = "/debug/allocs"
	RouteProfileBlock        = "/debug/block"
	RouteProfileCmdline      = "/debug/cmdline"
	RouteProfileGoroutine    = "/debug/goroutine"
	RouteProfileHeap         = "/debug/heap"
	RouteProfileMutex        = "/debug/mutex"
	RouteProfileProfile      = "/debug/profile"
	RouteProfileThreadcreate = "/debug/threadcreate"
	RouteProfileSymbol       = "/debug/symbol"
	RouteProfileTrace        = "/debug/trace"
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
		newRoute(http.MethodGet, RouteAsset, app.handleAssets()),

		newRoute(http.MethodGet, RouteHome, app.handleHome()),
		newRoute(http.MethodPost, RouteTodoCreate, app.handleTodoCreate()),
		newRoute(http.MethodGet, RouteTodoEdit, app.handleTodoEdit()),
		newRoute(http.MethodPost, RouteTodoUpdate, app.handleTodoUpdate()),
	}

	{
		observabilityRoutes := []route{
			newRoute(http.MethodGet, RouteHealth, app.handleHealth()),
			newRoute(http.MethodGet, RouteProfileBaseRoute, app.handleIndex()),
			newRoute(http.MethodGet, RouteProfileAllocs, app.handleAllocs()),
			newRoute(http.MethodGet, RouteProfileBlock, app.handleBlock()),
			newRoute(http.MethodGet, RouteProfileCmdline, app.handleCmdline()),
			newRoute(http.MethodGet, RouteProfileGoroutine, app.handleGoroutine()),
			newRoute(http.MethodGet, RouteProfileHeap, app.handleHeap()),
			newRoute(http.MethodGet, RouteProfileMutex, app.handleMutex()),
			newRoute(http.MethodGet, RouteProfileProfile, app.handleProfile()),
			newRoute(http.MethodGet, RouteProfileThreadcreate, app.handleThreadcreate()),
			newRoute(http.MethodGet, RouteProfileSymbol, app.handleSymbol()),
			newRoute(http.MethodGet, RouteProfileTrace, app.handleTrace()),
		}

		routes = append(routes, observabilityRoutes...)
	}

	routerEntry := app.routing(routes)

	return middleware.RequestLogger(
		middleware.RequestPath(
			middleware.RequestID(
				middleware.Recovery(
					middleware.RequestDuration(
						routerEntry, RouteAssetBase,
					),
				),
			), RouteAssetBase,
		),
	)
}

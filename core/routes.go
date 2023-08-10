package core

import (
	"net/http"
	"regexp"

	"github.com/roscrl/light/core/middleware"
)

const (
	RouteAssetBase = "/assets"
	RouteAsset     = "/assets/(.*)"

	RouteHome       = "/"
	RouteTodoCreate = "/todo/create"

	RouteUp                  = "/health"
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

type contextKeyFields struct{}

func (s *Server) routes() http.Handler {
	newRoute := func(method, pattern string, handler http.HandlerFunc) route {
		return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
	}

	routes := []route{
		newRoute(http.MethodGet, RouteAsset, s.handleAssets()),

		newRoute(http.MethodGet, RouteHome, s.handleHome()),
		newRoute(http.MethodPost, RouteTodoCreate, s.handleTodoCreate()),

		newRoute(http.MethodGet, RouteUp, s.handleHealth()),
	}

	{
		pprofRoutes := []route{
			newRoute(http.MethodGet, RouteProfileBaseRoute, s.handleIndex()),
			newRoute(http.MethodGet, RouteProfileAllocs, s.handleAllocs()),
			newRoute(http.MethodGet, RouteProfileBlock, s.handleBlock()),
			newRoute(http.MethodGet, RouteProfileCmdline, s.handleCmdline()),
			newRoute(http.MethodGet, RouteProfileGoroutine, s.handleGoroutine()),
			newRoute(http.MethodGet, RouteProfileHeap, s.handleHeap()),
			newRoute(http.MethodGet, RouteProfileMutex, s.handleMutex()),
			newRoute(http.MethodGet, RouteProfileProfile, s.handleProfile()),
			newRoute(http.MethodGet, RouteProfileThreadcreate, s.handleThreadcreate()),
			newRoute(http.MethodGet, RouteProfileSymbol, s.handleSymbol()),
			newRoute(http.MethodGet, RouteProfileTrace, s.handleTrace()),
		}

		routes = append(routes, pprofRoutes...)
	}

	routerEntry := s.routing(routes)

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

func getField(r *http.Request, index int) string {
	fields := r.Context().Value(contextKeyFields{}).([]string)

	return fields[index]
}

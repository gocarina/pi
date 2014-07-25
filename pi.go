package pi

import (
	"net/http"
	"github.com/gorilla/pat"
	"bytes"
	"sort"
)


// Pi represents the core of the API toolkit.
type Pi struct {
	router *pat.Router
	routes routes
}

// New returns a new Pi.
func New() *Pi {
	return &Pi{
		router: pat.New(),
	}
}

// Route adds a route to the Pi router.
func (p *Pi) Route(routeURL string, childRoutes ...*route) *route {
	route := Route(routeURL, childRoutes...)
	p.routes = append(p.routes, route)
	return route
}

// ListenAndServe listens on the TCP network address srv.Addr and then calls
// Serve to handle requests on incoming connections. If srv.Addr is blank, ":http" is used.
func (p *Pi) ListenAndServe(addr string) error {
	sort.Sort(p.routes)
	for _, route := range p.routes {
		p.constructPath(route)
	}
	return http.ListenAndServe(addr, p)
}

// ServeHTTP .
func (p *Pi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

func wrapHandler(handler HandlerFunction, routeURL string, parentRoutes ...*route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		context := newRequestContext(w, r, routeURL)
		onError := func() {
			for _, parentRoute := range parentRoutes {
				parentRoute.Interceptors.onError(context)
			}
		}
		for _, parentRoute := range parentRoutes {
			if err = parentRoute.Interceptors.before(context); err != nil {
				onError()
				return
			}
		}
		if err := handler(context); err != nil {
			onError()
			return
		}
		for _, parentRoute := range parentRoutes {
			if err = parentRoute.Interceptors.after(context); err != nil {
				break
			}
		}
	}
}

func (p *Pi) constructPath(parentRoutes ...*route) {
	lastRoute := parentRoutes[len(parentRoutes) - 1]
	for _, childRoute := range lastRoute.ChildRoutes {
		p.constructPath(append(parentRoutes, childRoute)...)
	}
	routeURLBuffer := bytes.Buffer{}
	for _, route := range parentRoutes {
		routeURLBuffer.WriteString(route.RouteURL)
	}
	routeURL := routeURLBuffer.String()
	for method, handler := range lastRoute.Methods {
		p.router.Add(method, routeURL, wrapHandler(handler, routeURL, parentRoutes...))
	}
}

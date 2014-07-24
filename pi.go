package pi

import (
	"net/http"
	"github.com/gorilla/pat"
	"bytes"
)

// Pi represents the core of the API toolkit.
type Pi struct {
	router *pat.Router
	routes []*route
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
	for _, route := range p.routes {
		p.constructPath(route)
	}
	return http.ListenAndServe(addr, p)
}

// ServeHTTP .
func (p *Pi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

func wrapHandler(handler HandlerFunction, parentRoutes ...*route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		context := newRequestContext(w, r)
		onError := func() {
			for _, parentRoute := range parentRoutes {
				parentRoute.interceptors.onError(context)
			}
		}
		for _, parentRoute := range parentRoutes {
			if err = parentRoute.interceptors.before(context); err != nil {
				onError()
				return
			}
		}
		if err := handler(context); err != nil {
			onError()
			return
		}
		for _, parentRoute := range parentRoutes {
			if err = parentRoute.interceptors.after(context); err != nil {
				break
			}
		}
	}
}

func (p *Pi) constructPath(parentRoutes ...*route) {
	lastRoute := parentRoutes[len(parentRoutes) - 1]
	for _, childRoute := range lastRoute.childRoutes {
		p.constructPath(append(parentRoutes, childRoute)...)
	}
	routeURL := bytes.Buffer{}
	for _, route := range parentRoutes {
		routeURL.WriteString(route.routeURL)
	}
	for method, handler := range lastRoute.methods {
		p.router.Add(method, routeURL.String(), wrapHandler(handler, parentRoutes...))
	}
}

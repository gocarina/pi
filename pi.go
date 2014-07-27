package pi

import (
	"bytes"
	"github.com/gorilla/pat"
	"net/http"
	"sort"
	"fmt"
	"os"
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

// Router adds a route to the Pi router.
func (p *Pi) Router(routeURL string, childRoutes ...*route) *route {
	route := newRoute(routeURL, childRoutes...)
	p.routes = append(p.routes, route)
	return route
}

// Route adds a subroute to a router or router.
func (p *Pi) Route(routeURL string, childRoutes ...*route) *route {
	return newRoute(routeURL, childRoutes...)
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

// ServeHTTP servers a route in the HTTP server.
func (p *Pi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.router.ServeHTTP(w, r)
}

// wrapHandler wraps a route handler to run the interceptors and the handler.
func wrapHandler(handler HandlerFunction, routeURL string, parentRoutes ...*route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		context := newRequestContext(w, r, routeURL)
		errorInterceptors := func(err error) {
			errorsHandled := false
			for _, parentRoute := range parentRoutes {
				errorsHandled = errorsHandled || parentRoute.Interceptors.runErrorInterceptors(context, err)
			}
			if !errorsHandled {
				if piError, ok := err.(HTTPError); ok {
					http.Error(context.W, piError.Err.Error(), piError.StatusCode)
				}
			}
		}
		for _, parentRoute := range parentRoutes {
			if err := parentRoute.Interceptors.runBeforeInterceptors(context); err != nil {
				errorInterceptors(err)
				return
			}
		}
		if err := handler(context); err != nil {
			errorInterceptors(err)
			return
		}
		for _, parentRoute := range parentRoutes {
			if err := parentRoute.Interceptors.runAfterInterceptors(context); err != nil {
				fmt.Sprintln(os.Stderr, "after interceptor raised error:", err)
			}
		}
	}
}

// constructPath constructs the path to the specified route/sub-route.
func (p *Pi) constructPath(parentRoutes ...*route) {
	lastRoute := parentRoutes[len(parentRoutes)-1]
	for _, childRoute := range lastRoute.ChildRoutes {
		p.constructPath(append(parentRoutes, childRoute)...)
	}
	routeURLBuffer := bytes.Buffer{}
	for _, route := range parentRoutes {
		routeURLBuffer.WriteString(route.RouteURL)
	}
	routeURL := routeURLBuffer.String()
	if len(routeURL) > 2 && routeURL[0:2] == "//" {
		routeURL = routeURL[1:]
	}
	for method, handler := range lastRoute.Methods {
		p.router.Add(method, routeURL, wrapHandler(handler, routeURL, parentRoutes...))
	}
}

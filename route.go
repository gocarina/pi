package pi

import "strings"

type route struct {
	RouteURL     string
	ChildRoutes  routes
	Methods      map[string]HandlerFunction
	Interceptors interceptors
}

type routes []*route

// Route adds a new sub route to a parent Route.
// See pi.Route to add root Routes.
func Route(RouteURL string, ChildRoutes ...*route) *route {
	return &route{
		RouteURL:    RouteURL,
		ChildRoutes: ChildRoutes,
		Methods:     make(map[string]HandlerFunction),
	}
}

// Before registers an HandlerFunction to be called before the request is handled.
func (r *route) Before(handler HandlerFunction) *route {
	r.Interceptors.addBefore(handler)
	return r
}

// After registers an HandlerFunction to be called after the request has been handled.
func (r *route) After(handler HandlerFunction) *route {
	r.Interceptors.addAfter(handler)
	return r
}

// OnError registers an HandlerFunction to be called when an error occurs in the main handler or any Before interceptor.
func (r *route) OnError(handler HandlerFunction) *route {
	r.Interceptors.addOnError(handler)
	return r
}

func (r *route) Any(handlerFunc HandlerFunction) *route {
	r.Get(handlerFunc)
	r.Post(handlerFunc)
	r.Put(handlerFunc)
	r.Delete(handlerFunc)
	r.Patch(handlerFunc)
	r.Options(handlerFunc)
	return r
}

// Get registers an HandlerFunction to handle GET requests.
func (r *route) Get(handlerFunc HandlerFunction) *route {
	r.Methods["GET"] = handlerFunc
	return r
}

// Post registers an HandlerFunction to handle POST requests.
func (r *route) Post(handlerFunc HandlerFunction) *route {
	r.Methods["POST"] = handlerFunc
	return r
}

// Put registers an HandlerFunction to handle PUT requests.
func (r *route) Put(handlerFunc HandlerFunction) *route {
	r.Methods["PUT"] = handlerFunc
	return r
}

// Delete registers an HandlerFunction to handle DELETE requests.
func (r *route) Delete(handlerFunc HandlerFunction) *route {
	r.Methods["DELETE"] = handlerFunc
	return r
}

// Patch registers an HandlerFunction to handle PATCH requests.
func (r *route) Patch(handlerFunc HandlerFunction) *route {
	r.Methods["PATCH"] = handlerFunc
	return r
}

// Options registers an HandlerFunction to handle OPTIONS requests.
func (r *route) Options(handlerFunc HandlerFunction) *route {
	r.Methods["OPTIONS"] = handlerFunc
	return r
}

func (r routes) Len() int {
	return len(r)
}

func (r routes) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r routes) Less(i, j int) bool {
	return strings.Count(r[i].RouteURL, "/") > strings.Count(r[j].RouteURL, "/")
}

package pi

import "strings"

// route represents an API route.
// For example: /user/get/{id}
type route struct {
	RouteURL     string
	ChildRoutes  routes
	Methods      map[string]HandlerFunction
	Interceptors interceptors
}

type routes []*route

// newRoute returns a new route.
func newRoute(RouteURL string, ChildRoutes ...*route) *route {
	return &route{
		RouteURL:    RouteURL,
		ChildRoutes: ChildRoutes,
		Methods:     make(map[string]HandlerFunction),
	}
}

// Before registers an interceptor to be called before the request is handled.
func (r *route) Before(b beforeInterceptor) *route {
	r.Interceptors.addBefore(b)
	return r
}

// After registers an interceptor to be called after the request has been handled.
func (r *route) After(a afterInterceptor) *route {
	r.Interceptors.addAfter(a)
	return r
}

// Error registers an interceptor to be called when an error occurs in the request handler or in any Before interceptor.
func (r *route) Error(e errorInterceptor) *route {
	r.Interceptors.addError(e)
	return r
}

// Intercept registers an interceptor to be called before, after or on error.
func (r *route) Intercept(ci interface {}) *route {
	r.Interceptors.addInterceptor(ci)
	return r
}

// Any registers an HandlerFunction to handle any requests.
func (r *route) Any(handlerFunc HandlerFunction) *route {
	r.Get(handlerFunc)
	r.Post(handlerFunc)
	r.Put(handlerFunc)
	r.Delete(handlerFunc)
	r.Patch(handlerFunc)
	r.Options(handlerFunc)
	r.Head(handlerFunc)
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

// Head registers an HandlerFunction to handle HEAD requests.
func (r *route) Head(handlerFunc HandlerFunction) *route {
	r.Methods["HEAD"] = handlerFunc
	return r
}

// Custom registers an HandlerFunction to handle custom requests.
func (r *route) Custom(method string, handlerFunc HandlerFunction) *route {
	r.Methods[method] = handlerFunc
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

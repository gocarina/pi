package pi

import "strings"

// Route represents an API Route.
// For example: /user/get/{id}
type Route struct {
	RouteURL     string
	ChildRoutes  routes
	Methods      map[string]HandlerFunction
	Interceptors interceptors
}

type routes []*Route

// newRoute returns a new Route.
func newRoute(RouteURL string, ChildRoutes ...*Route) *Route {
	return &Route{
		RouteURL:    RouteURL,
		ChildRoutes: ChildRoutes,
		Methods:     make(map[string]HandlerFunction),
	}
}

// interceptorHelper is an helper to add single funcs to the interceptor list.
type interceptorHelper struct {
	BeforeFunc HandlerFunction
	AfterFunc  HandlerFunction
	ErrorFunc  func(error) HandlerFunction
}

func (helper *interceptorHelper) Before() HandlerFunction {
	return helper.BeforeFunc
}

func (helper *interceptorHelper) After() HandlerFunction {
	return helper.AfterFunc
}

func (helper *interceptorHelper) Error(err error) HandlerFunction {
	return helper.ErrorFunc(err)
}

// Before registers an interceptor to be called before the request is handled.
func (r *Route) Before(b beforeInterceptor) *Route {
	r.Interceptors.addBefore(b)
	return r
}

// BeforeFunc add an an Before handler to the interceptor.
func (r *Route) BeforeFunc(handler HandlerFunction) *Route {
	helper := &interceptorHelper{
		BeforeFunc: handler,
	}
	r.Interceptors.addBefore(helper)
	return r
}

// After registers an interceptor to be called after the request has been handled.
func (r *Route) After(a afterInterceptor) *Route {
	r.Interceptors.addAfter(a)
	return r
}

// AfterFunc add an an After interceptor.
func (r *Route) AfterFunc(handler HandlerFunction) *Route {
	helper := &interceptorHelper{
		AfterFunc: handler,
	}
	r.Interceptors.addAfter(helper)
	return r
}

// Error registers an interceptor to be called when an error occurs in the request handler or in any Before interceptor.
func (r *Route) Error(e errorInterceptor) *Route {
	r.Interceptors.addError(e)
	return r
}

// ErrorFunc add an an Error interceptor.
func (r *Route) ErrorFunc(handler func (error) HandlerFunction) *Route {
	helper := &interceptorHelper{
		ErrorFunc: handler,
	}
	r.Interceptors.addError(helper)
	return r
}

// Intercept registers an interceptor to be called before, after or on error.
func (r *Route) Intercept(ci interface{}) *Route {
	r.Interceptors.addInterceptor(ci)
	return r
}

// Any registers an HandlerFunction to handle any requests.
func (r *Route) Any(handlerFunc HandlerFunction) *Route {
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
func (r *Route) Get(handlerFunc HandlerFunction) *Route {
	r.Methods["GET"] = handlerFunc
	return r
}

// Post registers an HandlerFunction to handle POST requests.
func (r *Route) Post(handlerFunc HandlerFunction) *Route {
	r.Methods["POST"] = handlerFunc
	return r
}

// Put registers an HandlerFunction to handle PUT requests.
func (r *Route) Put(handlerFunc HandlerFunction) *Route {
	r.Methods["PUT"] = handlerFunc
	return r
}

// Delete registers an HandlerFunction to handle DELETE requests.
func (r *Route) Delete(handlerFunc HandlerFunction) *Route {
	r.Methods["DELETE"] = handlerFunc
	return r
}

// Patch registers an HandlerFunction to handle PATCH requests.
func (r *Route) Patch(handlerFunc HandlerFunction) *Route {
	r.Methods["PATCH"] = handlerFunc
	return r
}

// Options registers an HandlerFunction to handle OPTIONS requests.
func (r *Route) Options(handlerFunc HandlerFunction) *Route {
	r.Methods["OPTIONS"] = handlerFunc
	return r
}

// Head registers an HandlerFunction to handle HEAD requests.
func (r *Route) Head(handlerFunc HandlerFunction) *Route {
	r.Methods["HEAD"] = handlerFunc
	return r
}

// Custom registers an HandlerFunction to handle custom requests.
func (r *Route) Custom(method string, handlerFunc HandlerFunction) *Route {
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

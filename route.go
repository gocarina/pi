package pi

type route struct {
	routeURL     string
	childRoutes  []*route
	methods      map[string]HandlerFunction
	interceptors requestInterceptors
}

// Route adds a new sub route to a parent Route.
// See pi.Route to add root Routes.
func Route(routeURL string, childRoutes ...*route) *route {
	return &route{
		routeURL:    routeURL,
		childRoutes: childRoutes,
		methods:     make(map[string]HandlerFunction),
	}
}

// Before registers an HandlerFunction to be called before the request is handled.
func (r *route) Before(handler HandlerFunction) *route {
	r.interceptors.addBefore(handler)
	return r
}

// After registers an HandlerFunction to be called after the request has been handled.
func (r *route) After(handler HandlerFunction) *route {
	r.interceptors.addAfter(handler)
	return r
}

// OnError registers an HandlerFunction to be called when an error occurs in the main handler or any Before interceptor.
func (r *route) OnError(handler HandlerFunction) *route {
	r.interceptors.addOnError(handler)
	return r
}

// Get registers an HandlerFunction to handle GET requests.
func (r *route) Get(handlerFunc HandlerFunction) *route {
	r.methods["GET"] = handlerFunc
	return r
}

// Post registers an HandlerFunction to handle POST requests.
func (r *route) Post(handlerFunc HandlerFunction) *route {
	r.methods["POST"] = handlerFunc
	return r
}

// Put registers an HandlerFunction to handle PUT requests.
func (r *route) Put(handlerFunc HandlerFunction) *route {
	r.methods["PUT"] = handlerFunc
	return r
}

// Delete registers an HandlerFunction to handle DELETE requests.
func (r *route) Delete(handlerFunc HandlerFunction) *route {
	r.methods["DELETE"] = handlerFunc
	return r
}

// Patch registers an HandlerFunction to handle PATCH requests.
func (r *route) Patch(handlerFunc HandlerFunction) *route {
	r.methods["PATCH"] = handlerFunc
	return r
}

// Options registers an HandlerFunction to handle OPTIONS requests.
func (r *route) Options(handlerFunc HandlerFunction) *route {
	r.methods["OPTIONS"] = handlerFunc
	return r
}

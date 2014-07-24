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

func (r *route) Before(handler HandlerFunction) *route {
	r.interceptors.addBefore(handler)
	return r
}

func (r *route) After(handler HandlerFunction) *route {
	r.interceptors.addAfter(handler)
	return r
}

func (r *route) OnError(handler HandlerFunction) *route {
	r.interceptors.addOnError(handler)
	return r
}

func (r *route) Get(handlerFunc HandlerFunction) *route {
	r.methods["GET"] = handlerFunc
	return r
}

func (r *route) Post(handlerFunc HandlerFunction) *route {
	r.methods["POST"] = handlerFunc
	return r
}
func (r *route) Put(handlerFunc HandlerFunction) *route {
	r.methods["PUT"] = handlerFunc
	return r
}
func (r *route) Delete(handlerFunc HandlerFunction) *route {
	r.methods["DELETE"] = handlerFunc
	return r
}

func (r *route) Patch(handlerFunc HandlerFunction) *route {
	r.methods["PATCH"] = handlerFunc
	return r
}

func (r *route) Options(handlerFunc HandlerFunction) *route {
	r.methods["OPTIONS"] = handlerFunc
	return r
}

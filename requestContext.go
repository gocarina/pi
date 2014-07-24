package pi

import "net/http"

// RequestContext represents the context of the HTTP request.
// It is shared across interceptors and handler.
type RequestContext struct {
	W        http.ResponseWriter
	R        *http.Request
	RouteURL string
	Data     map[interface{}]interface{}
}

func newRequestContext(w http.ResponseWriter, r *http.Request, routeURL string) *RequestContext {
	return &RequestContext{
		W:    w,
		R:    r,
		RouteURL: routeURL,
		Data: make(map[interface{}]interface{}),
	}
}

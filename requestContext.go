package pi

import "net/http"

// RequestContext represents the context of the HTTP request.
// It is shared across interceptors
type RequestContext struct {
	W    http.ResponseWriter
	R    *http.Request
	Data map[interface{}]interface{}
}

func newRequestContext(w http.ResponseWriter, r *http.Request) *RequestContext {
	return &RequestContext{
		W:w,
		R:r,
		Data: make(map[interface{}]interface{}),
	}
}

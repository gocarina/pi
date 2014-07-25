package pi

import (
	"net/http"
	"encoding/json"
)

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
		W:        w,
		R:        r,
		RouteURL: routeURL,
		Data:     make(map[interface{}]interface{}),
	}
}

// WriteJSON marshal the Object and write it.
func (c *RequestContext) WriteJSON(object interface{}) error {
	c.W.Header().Add("Content-Type", "application/json; charset=utf-8")
	if debugMode {
		js, err := json.MarshalIndent(object, "", "  ")
		if err != nil {
			return err
		}
		c.W.Write(js)
	} else {
		js, err := json.Marshal(object)
		if err != nil {
			return err
		}
		c.W.Write(js)
	}
	return nil
}

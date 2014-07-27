package pi

import (
	"fmt"
	"os"
)

// beforeInterceptor represents an interface with a Before() method called before the request handler.
type beforeInterceptor interface {
	Before() HandlerFunction
}

// afterInterceptor represents an interface with an After() method called after the request handler.
type afterInterceptor interface {
	After() HandlerFunction
}

// errorInterceptor represents an interface with an Error() method called when an error occurred in a Before interceptor, or in the request handler.
type errorInterceptor interface {
	Error(error) HandlerFunction
}

// interceptors gathers Before, After and Error interceptors.
type interceptors struct {
	Before []beforeInterceptor
	After  []afterInterceptor
	Error  []errorInterceptor
}

// addBefore appends a Before interceptor.
func (i *interceptors) addBefore(b beforeInterceptor) {
	i.Before = append(i.Before, b)
}

// addAfter appends an After interceptor.
func (i *interceptors) addAfter(a afterInterceptor) {
	i.After = append(i.After, a)
}

// addError appends an Error interceptor.
func (i *interceptors) addError(e errorInterceptor) {
	i.Error = append(i.Error, e)
}

// runBeforeInterceptors runs all the Before interceptors, breaking if an error is thrown.
func (i *interceptors) runBeforeInterceptors(c *RequestContext) error {
	for _, b := range i.Before {
		if err := b.Before()(c); err != nil {
			return err
		}
	}
	return nil
}

// runAfterInterceptors runs all the After interceptors, ignoring if an error is thrown.
func (i *interceptors) runAfterInterceptors(c *RequestContext) error {
	for _, a := range i.After {
		if err := a.After()(c); err != nil {
			return err
		}
	}
	return nil
}

// runAfterInterceptors runs all the Error interceptors, ignoring if an error is thrown.
func (i *interceptors) runErrorInterceptors(c *RequestContext, err error) bool {
	for _, e := range i.Error {
		if err := e.Error(err)(c); err != nil {
			fmt.Sprintln(os.Stderr, "error interceptor raised error:", err)
		}
	}
	return len(i.Error) > 0
}

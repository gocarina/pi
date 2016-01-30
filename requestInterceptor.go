package pi

import (
	"fmt"
	"os"
)

// interceptors gathers Before, After and Error interceptors.
type interceptors struct {
	Before     []HandlerFunction
	After      []HandlerFunction
	AfterAsync []HandlerFunction
	Recoverers []RecovererFunction
	Error      []HandlerErrorFunction
}

// addBefore appends a Before interceptor.
func (i *interceptors) addBefore(handler HandlerFunction) {
	i.Before = append(i.Before, handler)
}

// addAfter appends an After interceptor.
func (i *interceptors) addAfter(handler HandlerFunction) {
	i.After = append(i.After, handler)
}

// addAfterAsync appends an AfterAsync interceptor.
func (i *interceptors) addAfterAsync(handler HandlerFunction) {
	i.AfterAsync = append(i.AfterAsync, handler)
}

// addRecoverer appends a Recoverer interceptor.
func (i *interceptors) addRecoverer(recoverer RecovererFunction) {
	i.Recoverers = append(i.Recoverers, recoverer)
}

// addError appends an Error interceptor.
func (i *interceptors) addError(handler HandlerErrorFunction) {
	i.Error = append(i.Error, handler)
}

// runBeforeInterceptors runs all the Before interceptors, breaking if an error is thrown.
func (i *interceptors) runBeforeInterceptors(c *RequestContext) error {
	for _, b := range i.Before {
		if err := b(c); err != nil {
			return err
		}
	}
	return nil
}

// runAfterInterceptors runs all the After interceptors, ignoring if an error is thrown.
func (i *interceptors) runAfterInterceptors(c *RequestContext) error {
	for _, a := range i.After {
		if err := a(c); err != nil {
			return err
		}
	}
	return nil
}

// runAfterAsyncInterceptors runs all the AfterAsync interceptors, ignoring every errors.
func (i *interceptors) runAfterAsyncInterceptors(c *RequestContext) {
	for _, as := range i.AfterAsync {
		go as(c)
	}
}

// runRecovererInterceptors runs all the Recoverer interceptors.
func (i *interceptors) runRecovererInterceptors(c *RequestContext, recoverValue interface{}) {
	for _, r := range i.Recoverers {
		r(c, recoverValue)
	}
}

// runAfterInterceptors runs all the Error interceptors, ignoring if an error is thrown.
func (i *interceptors) runErrorInterceptors(c *RequestContext, err error) (returnError error) {
	for _, e := range i.Error {
		if err := e(c, err); err != nil {
			fmt.Sprintln(os.Stderr, "error interceptor raised error:", err)
			returnError = err
		}
	}
	return
}

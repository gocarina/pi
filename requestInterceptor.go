package pi

const (
	undef = iota
	constBefore
	constAfter
	constOnError
)

// interceptor represents a function to be called before or after a route, or on error.
type interceptor struct {
	Type    int
	Handler HandlerFunction
}

type interceptors []*interceptor

// newBefore returns a new Before interceptor.
func newBefore(handler HandlerFunction) *interceptor {
	return &interceptor{
		Type:    constBefore,
		Handler: handler,
	}
}

// newAfter returns a new After interceptor.
func newAfter(handler HandlerFunction) *interceptor {
	return &interceptor{
		Type:    constAfter,
		Handler: handler,
	}
}

// newAfter returns a new OnError interceptor.
func newOnError(handler HandlerFunction) *interceptor {
	return &interceptor{
		Type:    constOnError,
		Handler: handler,
	}
}

// addBefore appends a Before interceptor to the interceptors.
func (r *interceptors) addBefore(handler HandlerFunction) {
	*r = append(*r, newBefore(handler))
}

// addAfter appends an After interceptor to the interceptors.
func (r *interceptors) addAfter(handler HandlerFunction) {
	*r = append(*r, newAfter(handler))
}

// addOnError appends an OnError interceptor to the interceptors.
func (r *interceptors) addOnError(handler HandlerFunction) {
	*r = append(*r, newOnError(handler))
}

// before executes all Before interceptors
func (r *interceptors) before(c *RequestContext) error {
	for _, interceptor := range *r {
		if interceptor.Type == constBefore {
			if err := interceptor.Handler(c); err != nil {
				return err
			}
		}
	}
	return nil
}

// after executes all After interceptors
func (r *interceptors) after(c *RequestContext) error {
	for _, interceptor := range *r {
		if interceptor.Type == constAfter {
			if err := interceptor.Handler(c); err != nil {
				return err
			}
		}
	}
	return nil
}

// onError executes all OnError interceptors
func (r *interceptors) onError(c *RequestContext) error {
	for _, interceptor := range *r {
		if interceptor.Type == constOnError {
			if err := interceptor.Handler(c); err != nil {
				return err
			}
		}
	}
	return nil
}

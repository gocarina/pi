package pi

const (
	undef = iota
	constBefore
	constAfter
	constOnError
)

type interceptor struct {
	Type    int
	Handler HandlerFunction
}

type interceptors []*interceptor

func newBefore(handler HandlerFunction) *interceptor {
	return &interceptor{
		Type:    constBefore,
		Handler: handler,
	}
}

func newAfter(handler HandlerFunction) *interceptor {
	return &interceptor{
		Type:    constAfter,
		Handler: handler,
	}
}

func newOnError(handler HandlerFunction) *interceptor {
	return &interceptor{
		Type:    constOnError,
		Handler: handler,
	}
}

func (r *interceptors) addBefore(handler HandlerFunction) {
	*r = append(*r, newBefore(handler))
}

func (r *interceptors) addAfter(handler HandlerFunction) {
	*r = append(*r, newAfter(handler))
}

func (r *interceptors) addOnError(handler HandlerFunction) {
	*r = append(*r, newOnError(handler))
}

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

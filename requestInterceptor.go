package pi

const (
	undef = iota
	constBefore
	constAfter
	constOnError
)

type requestInterceptor struct {
	Type    int
	Handler HandlerFunction
}

type requestInterceptors []*requestInterceptor

func newBefore(handler HandlerFunction) *requestInterceptor {
	return &requestInterceptor{
		Type:    constBefore,
		Handler: handler,
	}
}

func newAfter(handler HandlerFunction) *requestInterceptor {
	return &requestInterceptor{
		Type:    constAfter,
		Handler: handler,
	}
}

func newOnError(handler HandlerFunction) *requestInterceptor {
	return &requestInterceptor{
		Type:    constOnError,
		Handler: handler,
	}
}

func (r *requestInterceptors) addBefore(handler HandlerFunction) {
	*r = append(*r, newBefore(handler))
}

func (r *requestInterceptors) addAfter(handler HandlerFunction) {
	*r = append(*r, newAfter(handler))
}

func (r *requestInterceptors) addOnError(handler HandlerFunction) {
	*r = append(*r, newOnError(handler))
}

func (r *requestInterceptors) before(c *RequestContext) error {
	for _, rI := range *r {
		if rI.Type == constBefore {
			if err := rI.Handler(c); err != nil {
				return err
			}
		}
	}
	return nil
}
func (r *requestInterceptors) after(c *RequestContext) error {
	for _, rI := range *r {
		if rI.Type == constAfter {
			if err := rI.Handler(c); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *requestInterceptors) onError(c *RequestContext) error {
	for _, rI := range *r {
		if rI.Type == constOnError {
			if err := rI.Handler(c); err != nil {
				return err
			}
		}
	}
	return nil
}

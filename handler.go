package pi

// HandlerFunction represents.
type HandlerFunction func(*RequestContext) error

// SimpleHandler returns.
func SimpleHandler(f func (*RequestContext) string) HandlerFunction {
	return func(c *RequestContext) error {
		c.W.Write([]byte(f(c)))
		return nil
	}
}

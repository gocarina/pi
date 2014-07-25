package pi

import "testing"

func test1Handler(c *RequestContext) error {
	c.W.Write([]byte("test1"))
	return nil
}

func test2Handler(c *RequestContext) error {
	c.W.Write([]byte("user"))
	return nil
}

func TestPi(t *testing.T) {
	p := New()

	p.Route("/test",
		p.Route("/user").Delete(test2Handler).Get(test1Handler)).
		Get(test1Handler)
	p.ListenAndServe(":9001")
}

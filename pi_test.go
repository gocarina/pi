package pi

import (
	"testing"
)

func rootHandler(c *RequestContext) error {
	c.WriteString("/")
	return nil
}

func userHandler(c *RequestContext) error {
	c.WriteString("/user")
	return nil
}

func userIDHandler(c *RequestContext) error {
	c.WriteString("/user/" + c.GetRouteVariable("id"))
	return nil
}

func TestPi(t *testing.T) {
	p := New()
	router := p.Router("/",
		p.Route("/user",
			p.Route("/{id}").Get(userIDHandler).Post(userIDHandler).Put(userIDHandler).Delete(userIDHandler),
		).Get(userHandler),
	).Get(rootHandler)


	if router.ChildRoutes[0].RouteURL != "/user" {
		t.Fatal("Router SubRouting failed")
	}
	if router.ChildRoutes[0].ChildRoutes[0].RouteURL != "/{id}" {
		t.Fatal("Router SubRoute SubRouting failed")
	}
	if router.ChildRoutes[0].ChildRoutes[0].Methods["GET"] == nil &&
		router.ChildRoutes[0].ChildRoutes[0].Methods["POST"] == nil &&
		router.ChildRoutes[0].ChildRoutes[0].Methods["PUT"] == nil &&
		router.ChildRoutes[0].ChildRoutes[0].Methods["DELETE"] == nil {
		t.Fatal("Router SubRoute SubRoute method missing")
	}

	panic(p.ListenAndServe(":9001"))
}

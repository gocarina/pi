package pi

import (
	"log"
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

func test1(c *RequestContext) error {
	return c.WriteString("test1")
}

func test2(c *RequestContext) error {
	return c.WriteString("test2")
}

func test3(c *RequestContext) error {
	return c.WriteString("test3")
}

func before1(c *RequestContext) error {
	log.Println("Before1")
	return nil
}

func before3(c *RequestContext) error {
	log.Println("Before3")
	return nil
}

func TestPi(t *testing.T) {
	p := New()
	router := p.Router("/",
		p.Route("/user",
			p.Route("/{id}",
				p.Route("/test1").Get(test1).Before(before1),
				p.Route("/test2").Get(test2),
				p.Route("/test3").Get(test3).Before(before3),
			).Get(userIDHandler).Post(userIDHandler).Put(userIDHandler).Delete(userIDHandler),
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

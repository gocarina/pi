/*
Package pi is a powerful toolkit to fasten the writing of web services in Go.
It has been written in top of the router from the pat toolkit (http://gorillatoolkit.org/pkg/pat).


Simple example:

		import "github.com/gocarina/pi"

		func main() {
			p := pi.New()
			p.Router("/",

				p.Route("/products",

						// "/products/{id}"
						p.Route("/{id}").
						Delete(DeleteProductHandler).
						Put(EditProductHandler).
						BeforeFunc(AuthorizeUser)).

				// "/products"
				Get(GetProductsHandler).
				Post(AddProductHandler))

				p.ListenAndServe(":8080")
		}


This example registers two different routes :

- "/products" with two methods: GetProductsHandler (GET) and AddProductHandler (POST)

- "/products/{id}" with also two methods: DeleteProductHandler (GET) and EditProductHandler (PUT)

There is also an interceptor Before in the route "/products/{id}", that means that the function AuthorizeUser is called before calling either
DeleteProductHandler or EditProductHandler.

If there was interceptors on the route "/products", it would also apply for the child routes.
There is 3 kind of interceptor:


- Before: called before the request is handled by the handler. If there is an error in a Before interceptor, the request flow is stopped and
Error interceptors are called.

- After: called after the request has been handled by the handled. Errors are ignored (print to the error output).

- Error: called when an error occurs in any Before interceptor or in the request handler. Errors are completely ignored.


Interceptors and Handlers are both HandlerFunction.
Here is a typical example of HandlerFunction:

		func SomeHandlerOrInterceptor(c *pi.RequestContext) (err error) {
			// Do something
			if err != nil {
				return pi.NewError(http.StatusInternalServerError, err) // Output: { "errorCode": 500, "errorMessage": "Message of the error" }
			}
			return c.WriteJSON(pi.J{"status": "OK"}) // Output: { "status": "OK" }
		}

The RequestContext has some useful methods, please, see below for more information.

If you have any problems or questions, create an issue to our github page project: https://github.com/gocarina/pi

*/
package pi

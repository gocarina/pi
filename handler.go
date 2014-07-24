package pi

import "net/http"

// HandlerFunction represents the function called when the specified route is reached.
type HandlerFunction func(*RequestContext) error

// FileServeHandler registers the specified folder to be served on the specified route.
func FileServeHandler(path string) HandlerFunction {
	return func(c *RequestContext) error {
		file := c.R.URL.Query().Get(":file")
		if file != "" {
			http.ServeFile(c.W, c.R, path)
		} else {
			http.ServeFile(c.W, c.R, path+"/"+file)
		}
		return nil
	}
}

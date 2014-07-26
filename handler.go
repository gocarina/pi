package pi

import (
	"net/http"
)

// The HandlerFunction type is an adapter to allow the use of ordinary functions as route handlers.
type HandlerFunction func(*RequestContext) error

// ServeFileHandler replies to the request with the contents of the named file or directory.
// For example:
// p := New()
// p.Router("/files").Get(ServeFileHandler("/tmp"))
func ServeFileHandler(path string, allowBrowsing bool) HandlerFunction {
	return func(c *RequestContext) error {
		if allowBrowsing {
			fullPath := c.R.URL.String()
			if len(fullPath) > len(c.RouteURL) {
				path += fullPath[len(c.RouteURL):]
			}
		}
		http.ServeFile(c.W, c.R, path)
		return nil
	}
}

// ErrorHandler returns the status code string representation and set the status code as specified.
// See http.StatusText.
func ErrorHandler(statusCode int) HandlerFunction {
	return func(c *RequestContext) error {
		c.SetStatusCode(statusCode)
		c.WriteString(http.StatusText(statusCode))
		return nil
	}
}

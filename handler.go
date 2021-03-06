package pi

import (
	"fmt"
	"net/http"
	"os"
)

// The HandlerFunction type is an adapter to allow the use of ordinary functions as route handlers.
type HandlerFunction func(*RequestContext) error

// RecovererFunction type is a function containing the RequestContext of the request that panicked
// and the value recovered.
type RecovererFunction func(*RequestContext, interface{})

// HandlerErrorFunction type is the type used by error interceptors.
type HandlerErrorFunction func(*RequestContext, error) error

// ServeFileHandler replies to the request with the contents of the named file or directory.
// For example:
// p := New()
// p.Router("/files").Get(ServeFileHandler("/tmp"))
func ServeFileHandler(path string, allowBrowsing bool) HandlerFunction {
	if debugMode {
		if allowBrowsing {
			fmt.Fprintln(os.Stderr, "ServeFileHandler is vulnerable to Directory traversal attack when using allowBrowsing, use with caution!")
		}
	}
	return func(c *RequestContext) error {
		if allowBrowsing {
			path += c.GetRouteExtraPath()
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

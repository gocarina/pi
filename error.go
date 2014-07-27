package pi

import "fmt"

// HTTPError represents a HTTP Error.
type HTTPError struct {
	StatusCode int `json:"statusCode"`
	Err        error
}

// Error implements the Error interface.
func (h HTTPError) Error() string {
	return fmt.Sprintf("[%d] %s", int64(h.StatusCode), h.Err.Error())
}

// NewError returns a new HTTPError.
func NewError(statusCode int, err error) HTTPError {
	return HTTPError{
		StatusCode: statusCode,
		Err:        err,
	}
}

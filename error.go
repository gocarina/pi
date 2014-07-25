package pi

import "fmt"

// HTTPError represents a HTTP Error.
type HTTPError struct {
	StatusCode int64 `json:"statusCode"`
	Err        error
}

func (h HTTPError) Error() string {
	return fmt.Sprintf("[%d] ", h.StatusCode) + h.Err.Error()
}

// NewError returns a new HTTPError.
func NewError(statusCode int64, err error) *HTTPError {
	return &HTTPError{
		StatusCode: statusCode,
		Err:        err,
	}
}

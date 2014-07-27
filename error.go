package pi

import (
	"fmt"
)

// HTTPError represents a HTTP Error.
type HTTPError struct {
	StatusCode    int `json:"statusCode"`
	Err           error
	ErrorTemplate string
	ContentType   string
}

// Error implements the Error interface.
func (h HTTPError) Error() string {
	return fmt.Sprintf(h.ErrorTemplate, int64(h.StatusCode), h.Err.Error())
}

// NewError returns a new HTTPError, and outputs it as JSON.
func NewError(statusCode int, err error) HTTPError {
	return HTTPError{
		StatusCode:    statusCode,
		Err:           err,
		ErrorTemplate: `{"errorCode": %d, "errorMessage": %s}`,
		ContentType:   "application/json; charset=UTF-8",
	}
}

// NewXMLError returns a new HTTPError, and outputs it as XML.
func NewXMLError(statusCode int, err error) HTTPError {
	return HTTPError{
		StatusCode:    statusCode,
		Err:           err,
		ErrorTemplate: `<error code="%d">%s</error>`,
		ContentType:   "application/xml; charset=UTF-8",
	}
}

package pi

import (
	"fmt"
	"strconv"
)

// HTTPError represents a HTTP Error.
type HTTPError interface {
	error
	ContentType() string
	StatusCode() int
}

type _JSONError struct {
	statusCode int
	err        string
	template   string
}

func (error _JSONError) Error() string {
	return fmt.Sprintf(error.template, int(error.statusCode), strconv.Quote(error.err))
}

func (error _JSONError) StatusCode() int {
	return error.statusCode
}

func (error _JSONError) ContentType() string {
	return "application/json; charset=UTF-8"
}

type _XMLError struct {
	statusCode int
	err        string
	template   string
}

func (error _XMLError) Error() string {
	return fmt.Sprintf(error.template, int(error.statusCode), error.err)
}

func (error _XMLError) ContentType() string {
	return "application/xml; charset=UTF-8"
}

func (error _XMLError) StatusCode() int {
	return error.statusCode
}

// NewError returns a new HTTPError, and outputs it as JSON.
func NewError(statusCode int, err error) HTTPError {
	return _JSONError{
		statusCode: statusCode,
		err:        err.Error(),
		template:   `{"errorCode": %d, "errorMessage": %s}`,
	}
}

// NewXMLError returns a new HTTPError, and outputs it as XML.
func NewXMLError(statusCode int, err error) HTTPError {
	return _XMLError{
		statusCode: statusCode,
		err:        err.Error(),
		template:   `<error code="%d">%s</error>`,
	}
}

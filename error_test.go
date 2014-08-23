package pi

import (
	"fmt"
	"strconv"
	"testing"
)

func TestPiError(t *testing.T) {
	//	p := New()
	//	p.Router("/").Get(func(c *RequestContext) error {
	//		return NewError(418, fmt.Errorf("I'm a teapot"))
	//	}).Post(func(c *RequestContext) error {
	//		return NewXMLError(418, fmt.Errorf("I'm a teapot"))
	//	})
	//	p.ListenAndServe(":8080")
}

type customError struct {
	statusCode int
	err        string
}

func (error customError) StatusCode() int {
	return error.statusCode
}

func (error customError) ContentType() string {
	return "application/json; charset=utf-8"
}

func (error customError) Error() string {
	return fmt.Sprintf(`{"errorCode": "%d", "errorMessage": %s, "customMessage": "%s"}`, error.statusCode, strconv.Quote(error.err), "yolo")
}

func TestCustomError(t *testing.T) {
	p := New()
	p.Router("/").Get(func(c *RequestContext) error {
		return customError{
			statusCode: 418,
			err:        "this is sparta",
		}
	})
	p.ListenAndServe(":8080")
}

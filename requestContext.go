package pi

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// RequestContext represents the context of the HTTP request.
// It is shared across interceptors and handler.
type RequestContext struct {
	W        http.ResponseWriter
	R        *http.Request
	RouteURL string
	Data     map[interface{}]interface{}
}

func newRequestContext(w http.ResponseWriter, r *http.Request, routeURL string) *RequestContext {
	return &RequestContext{
		W:        w,
		R:        r,
		RouteURL: routeURL,
		Data:     make(map[interface{}]interface{}),
	}
}

// WriteJSON marshal the object to JSON and writes it via the ResponseWriter.
func (c *RequestContext) WriteJSON(object interface{}) error {
	c.W.Header().Add("Content-Type", "application/json; charset=utf-8")
	if debugMode {
		output, err := json.MarshalIndent(object, "", "  ")
		if err != nil {
			return err
		}
		writeDebug("WriteJSON", c.R.RemoteAddr, string(output))
		c.W.Write(output)
	} else {
		output, err := json.Marshal(object)
		if err != nil {
			return err
		}
		c.W.Write(output)
	}
	return nil
}

// WriteXML marshal the object to XML and writes it via the ResponseWriter.
func (c *RequestContext) WriteXML(object interface{}) error {
	c.W.Header().Add("Content-Type", "application/xml; charset=utf-8")
	if debugMode {
		output, err := xml.MarshalIndent(object, "", "  ")
		if err != nil {
			return err
		}
		writeDebug("WriteXML", c.R.RemoteAddr, string(output))
		c.W.Write(output)
	} else {
		output, err := xml.Marshal(object)
		if err != nil {
			return err
		}
		c.W.Write(output)
	}
	return nil
}

// GetBody return the body as a ReadCloser. It is the client responsibility to close the body.
func (c *RequestContext) GetBody() io.ReadCloser { return c.R.Body }

// GetRawBody return the body as a byte array. The body is already close.
func (c *RequestContext) GetRawBody() ([]byte, error) {
	body := c.GetBody()
	defer body.Close()
	rawBody, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	if debugMode {
		writeDebug("GetRawBody", c.R.RemoteAddr, fmt.Sprintf("got %s", string(rawBody)))
	}
	return ioutil.ReadAll(body)
}

// GetJSONObject call json.Unmarshal by sending the reference of the given object.
// For example:
//		func GetUser(c *pi.RequestContext) error {
// 			user := &User{}
// 			if err := c.GetJSONObject(user); err != nil {
//				return pi.NewError(400, err)
//			}
//			// Do something with the user...
//			return nil
//		}
//
func (c *RequestContext) GetJSONObject(object interface{}) error {
	rawBody, err := c.GetRawBody()
	if err != nil {
		return err
	}
	return json.Unmarshal(rawBody, &object)
}

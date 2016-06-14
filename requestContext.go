package pi

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gocarina/formdata"
	"github.com/gorilla/schema"
	"html/template"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

var (
	decoderFormValues = schema.NewDecoder()

	// ErrNoFiles is the error when no files is present in the value provided when
	// calling RequestContext.GetFileHeaders.
	ErrNoFiles = fmt.Errorf("no files")

	// ErrContentTypeNotSupported is the error when the format of the Content-Type is not supported
	ErrContentTypeNotSupported = fmt.Errorf("format not supported")

	// ContentTypeJSON is the default MIME for JSON data.
	ContentTypeJSON = "application/json"

	// ContentTypeXML is the default MIME for XML data.
	ContentTypeXML = "application/xml"

	// ContentTypeClassicForm is the default MIME for Form Encoded data.
	ContentTypeClassicForm = "application/x-www-form-urlencoded"

	// ContentTypeMultipart is the default MIME for Webform.
	ContentTypeMultipart = "multipart/form-data"

	// ContentTypeText is the default MIME for Textual data.
	ContentTypeText = "text/plain"
)

// J is an helper to write JSON.
// Example:
//		c.WriteJSON(pi.J{"status": "OK"})
//
type J map[string]interface{}

// RequestContext represents the context of the HTTP request.
// It is shared across interceptors and handler.
type RequestContext struct {
	W        http.ResponseWriter
	R        *http.Request
	RouteURL string
	Data     map[interface{}]interface{}
}

// newRequestContext returns a new RequestContext.
func newRequestContext(w http.ResponseWriter, r *http.Request, routeURL string) *RequestContext {
	return &RequestContext{
		W:        w,
		R:        r,
		RouteURL: routeURL,
		Data:     make(map[interface{}]interface{}),
	}
}

// WriteString writes the specified strings to the ResponseWriter.
func (c *RequestContext) WriteString(strings ...string) error {
	for _, s := range strings {
		if _, err := c.W.Write([]byte(s)); err != nil {
			return err
		}
	}
	return nil
}

// WriteJSON marshal the object to JSON and writes it via the ResponseWriter.
func (c *RequestContext) WriteJSON(object interface{}) error {
	c.W.Header().Set("Content-Type", "application/json; charset=utf-8")
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
	c.W.Header().Set("Content-Type", "application/xml; charset=utf-8")
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

// WriteDefault writes the object to the caller according to the acceptable MIME in the Accept header value.
// If the MIME is not supported, it sends a 406 Not Acceptable request.
// text/plain uses the String method to be serialized.
// Mime supported for write:
//
//		application/json
//		application/xml
//		text/plain
//
// If no Accept header is present, it writes the object as JSON.
func (c *RequestContext) WriteDefault(object interface{}) error {
	acceptedContentType := c.GetAccepts()
	for _, contentType := range acceptedContentType {
		switch contentType {
		case ContentTypeText:
			return c.WriteString(fmt.Sprintf("%v", object))
		case ContentTypeJSON, "*/*":
			return c.WriteJSON(object)
		case ContentTypeXML:
			return c.WriteXML(object)
		}
	}
	return NewError(406, ErrContentTypeNotSupported)
}

// WriteReader copy the reader to the ResponseWriter.
func (c *RequestContext) WriteReader(reader io.Reader) error {
	_, err := io.Copy(c.W, reader)
	return err
}

// WriteTemplate writes the given template to the ResponseWriter.
func (c *RequestContext) WriteTemplateFile(filename string, data interface{}) error {
	tmplate, err := template.ParseFiles(filename)
	if err != nil {
		return err
	}
	return tmplate.Execute(c.W, data)
}

// SetStatusCode sets the response status code.
func (c *RequestContext) SetStatusCode(statusCode int) {
	c.W.WriteHeader(statusCode)
}

// GetBody return the body as a ReadCloser. It is the client responsibility to close the body.
func (c *RequestContext) GetBody() io.ReadCloser {
	return c.R.Body
}

// GetRawBody returns the body as a byte array, closing the body reader.
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
	return rawBody, nil
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

// GetXMLObject call xml.Unmarshal by sending the reference of the given object.
func (c *RequestContext) GetXMLObject(object interface{}) error {
	rawBody, err := c.GetRawBody()
	if err != nil {
		return err
	}
	return xml.Unmarshal(rawBody, &object)
}

// GetFormObject call gorilla/schema.Decode to maps the form values of the request into the object.
// For more information: http://www.gorillatoolkit.org/pkg/schema
func (c *RequestContext) GetFormObject(object interface{}) error {
	if err := c.R.ParseForm(); err != nil {
		return err
	}
	return decoderFormValues.Decode(object, c.R.PostForm)
}

// GetMultipartObject calls gocarina/formdata.Unmarshal to maps the multipart form values of the request into the object.
// It supports files through multipart.FileHeader.
func (c *RequestContext) GetMultipartObject(object interface{}) error {
	return formdata.Unmarshal(c.R, object)
}

// GetDefaultObject calls one of a GetX method according to the Content-Type of the request.
// Content-Types supported:
// 		application/json
//		application/xml
//		application/x-www-form-urlencoded
//		multipart/form-data
//
func (c *RequestContext) GetDefaultObject(object interface{}) error {
	switch c.GetContentType() {
	case ContentTypeJSON:
		return c.GetJSONObject(object)
	case ContentTypeXML:
		return c.GetXMLObject(object)
	case ContentTypeClassicForm:
		return c.GetFormObject(object)
	}
	if strings.SplitAfter(c.GetContentType(), ContentTypeMultipart)[0] == ContentTypeMultipart {
		return c.GetMultipartObject(object)
	}
	return ErrContentTypeNotSupported
}

// GetRouteExtraPath returns the extra path.
// For example:
// 		for route("/files"), "/files/home/user/.emacs" will return "/home/user/.emacs"
func (c *RequestContext) GetRouteExtraPath() (path string) {
	fullPath := c.R.URL.String()
	if len(fullPath) > len(c.RouteURL) {
		path = fullPath[len(c.RouteURL):]
	}
	return path
}

// GetURLParam returns an URL parameter.
// For example, given this URL:
//		/user?id=1234
//
//		fmt.Println(c.GetURLParam("id"))
//		// Outputs 1234
//
func (c *RequestContext) GetURLParam(param string) string {
	return c.R.FormValue(param)
}

// GetURLParams returns multiple URL parameter.
// For example, given this URL:
//		/user?c=1&c=2
//
//		fmt.Println(c.GetURLParams("c"))
//		// Outputs [1, 2]
//
func (c *RequestContext) GetURLParams(param string) []string {
	c.R.ParseForm()
	return c.R.Form[param]
}

// GetURLParamOrDefault returns an URL parameter or the defaultValue if the value is empty.
func (c *RequestContext) GetURLParamOrDefault(param, defaultValue string) string {
	value := c.GetURLParam(param)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetURLParamAsInt returns the URL parameter as an int. Ignores error.
func (c *RequestContext) GetURLParamAsInt(param string) int {
	v, _ := strconv.Atoi(c.GetURLParam(param))
	return v
}

// GetRouteVariable returns a route variable.
// For example:
//		getUserByID := func(c *RequestContext) error {
//			id := c.GetRouteVariable("id")
//			// Do something with the ID.
//			return nil
//		}
//
//		p := pi.New()
//		p.Route("/user/{id}").Get(getUserByID)
//		p.ListenAndServe(":8080")
//
func (c *RequestContext) GetRouteVariable(key string) string {
	return c.R.URL.Query().Get(":" + key)
}

// GetFileHeaders returns an array of FileHeader.
// For example:
//		getImages := func(c *RequestContext) error {
//			files, err := c.GetFileHeaders("files")
//			if err != nil {
//				return pi.NewError(400, err)
//			}
//			// Handle files
//			return nil
//		}
//
//		p := pi.New()
//		p.Route("/images").Post(getImages)
//		p.ListenAndServe(":8080")
//
func (c *RequestContext) GetFileHeaders(key string) ([]*multipart.FileHeader, error) {
	if err := c.R.ParseMultipartForm(32 << 20); err != nil {
		return nil, err
	}
	if c.R.MultipartForm != nil && c.R.MultipartForm.File[key] != nil {
		if debugMode {
			writeDebug("GetFileHeaders", c.R.RemoteAddr, fmt.Sprintf("got %d files from key %s", len(c.R.MultipartForm.File[key]), key))
		}
		return c.R.MultipartForm.File[key], nil
	}
	if debugMode {
		writeDebug("GetFileHeaders", c.R.RemoteAddr, fmt.Sprintf("got 0 files from key %s", key))
	}
	return nil, ErrNoFiles
}

// GetHeader returns the first value of the header of the request associated with the given key.
func (c *RequestContext) GetHeader(key string) string {
	return c.R.Header.Get(key)
}

// GetContentType returns the Content-Type of the request.
func (c *RequestContext) GetContentType() string {
	return c.GetHeader("Content-Type")
}

// GetAccepts returns the MIME types acceptable by the caller.
func (c *RequestContext) GetAccepts() []string {
	return c.GetHeaders("Accept")
}

// GetHeaderOrDefault returns the value of the header of the given defaultValue if the value is empty.
func (c *RequestContext) GetHeaderOrDefault(key, defaultValue string) string {
	value := c.GetHeader(key)
	if value == "" {
		return defaultValue
	}
	return key
}

// SetHeader sets to the response the header entries associated with key to
// the single element value.  It replaces any existing
// values associated with key.
func (c *RequestContext) SetHeader(key, value string) {
	c.W.Header().Set(key, value)
}

// AddHeader adds to the response the key, value pair to the header.
// It appends to any existing values associated with key.
func (c *RequestContext) AddHeader(key, value string) {
	c.W.Header().Add(key, value)
}

// DeleteHeader deletes from the response the values associated with key.
func (c *RequestContext) DeleteHeader(key string) {
	c.W.Header().Del(key)
}

// GetHeaders returns the array of values associated with the key from the request.
func (c *RequestContext) GetHeaders(key string) []string {
	return c.R.Header[key]
}

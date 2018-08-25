package app

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
)

// Context -
type Context struct {
	Config
}

// ConnectionString gets connection string for given key
func (c *Context) ConnectionString(key string) string {
	if val, ok := c.ConnectionStrings[key]; ok {
		return val
	}
	return ""
}

// AppSetting gets appSetting string for given key
func (c *Context) AppSetting(key string) string {
	if val, ok := c.AppSettings[key]; ok {
		return val
	}
	return ""
}

// ResponseData creates response success object
func (c *Context) ResponseData(data interface{}) Response {
	return Response{
		Success: true,
		Data:    data,
	}
}

// ResponseError creates response error object
func (c *Context) ResponseError(err error) Response {
	return Response{
		Success: false,
		Error:   err.Error(),
	}
}

// JSON marshals 'v' to JSON, automatically escaping HTML and setting the
// Content-Type as application/json.
func (c *Context) JSON(w http.ResponseWriter, status int, v interface{}) {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

// Plain writes a string to the response, setting the Content-Type as
// text/plain.
func (c *Context) Plain(w http.ResponseWriter, status int, v string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(v))
}

// XML marshals 'v' to JSON, setting the Content-Type as application/xml. It
// will automatically prepend a generic XML header (see encoding/xml.Header) if
// one is not found in the first 100 bytes of 'v'.
func (c *Context) XML(w http.ResponseWriter, status int, v interface{}) {
	b, err := xml.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(status)

	// Try to find <?xml header in first 100 bytes (just in case there're some XML comments).
	findHeaderUntil := len(b)
	if findHeaderUntil > 100 {
		findHeaderUntil = 100
	}
	if !bytes.Contains(b[:findHeaderUntil], []byte("<?xml")) {
		// No header found. Print it out first.
		w.Write([]byte(xml.Header))
	}

	w.Write(b)
}

// Data writes raw bytes to the response, setting the Content-Type as
// application/octet-stream.
func (c *Context) Data(w http.ResponseWriter, status int, v []byte) {
	w.Header().Set("Content-Type", "application/octet-stream")
	w.WriteHeader(status)
	w.Write(v)
}

// HTML writes a string to the response, setting the Content-Type as text/html.
func (c *Context) HTML(w http.ResponseWriter, status int, v string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(v))
}

// DecodeJSON decodes json data to object
func (c *Context) DecodeJSON(r io.Reader, v interface{}) error {
	defer io.Copy(ioutil.Discard, r)
	return json.NewDecoder(r).Decode(v)
}

// DecodeXML decodes XML data to object
func (c *Context) DecodeXML(r io.Reader, v interface{}) error {
	defer io.Copy(ioutil.Discard, r)
	return xml.NewDecoder(r).Decode(v)
}

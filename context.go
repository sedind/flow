package flow

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/sedind/flow/dbe"
	"github.com/sedind/flow/logger"
)

// Context -
type Context struct {
	Config
	DBConnections map[string]*dbe.Connection
	Logger        logger.Logger
}

// DefaultConnection gets default DB Connection
func (c *Context) DefaultConnection() *dbe.Connection {
	if c, ok := c.DBConnections[c.Config.DefaultConnection]; ok {
		return c
	}
	return nil
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

// Bind decodes a request body and binds it with v Object
func (c *Context) Bind(r *http.Request, v interface{}) error {
	ct := r.Header.Get("Content-Type")
	s := strings.TrimSpace(strings.Split(ct, ";")[0])
	switch s {
	case "application/json", "text/javascript":
		return c.DecodeJSON(r.Body, v)
	case "text/xml", "application/xml":
		return c.DecodeXML(r.Body, v)
	default:
		return errors.Errorf("Unsupported Content-Type: %s", s)
	}
}

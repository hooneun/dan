package dan

import (
	"encoding/json"
	"net/http"
)

type Context struct {
	W http.ResponseWriter
	R *http.Request
}

func (c *Context) JSON(statusCode int, data any) error {
	if c.Written() {
		return nil
	}

	c.W.Header().Set("Content-Type", "application/json")
	c.W.WriteHeader(statusCode)
	if data != nil {
		return json.NewEncoder(c.W).Encode(data)
	}

	return nil
}

func (c *Context) Error(statusCode int, message string) error {
	return c.JSON(statusCode, map[string]string{"error": message})
}

func (c *Context) Param(key string) string {
	return c.R.PathValue(key)
}

func (c *Context) Query(key string) string {
	return c.R.URL.Query().Get(key)
}

func (c *Context) DefaultQuery(key, defaultValue string) string {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func (c *Context) Form(key string) string {
	if err := c.R.ParseForm(); err != nil {
		return ""
	}

	return c.R.FormValue(key)
}

func (c *Context) DefaultForm(key, defaultValue string) string {
	value := c.Form(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func (c *Context) BindJSON(v any) error {
	defer c.R.Body.Close()
	return json.NewDecoder(c.R.Body).Decode(v)
}

func (c *Context) Written() bool {
	writer, ok := c.W.(interface{ Written() bool })
	return ok && writer.Written()
}

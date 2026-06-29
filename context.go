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

func (c *Context) BindJSON(v any) error {
	defer c.R.Body.Close()
	return json.NewDecoder(c.R.Body).Decode(v)
}

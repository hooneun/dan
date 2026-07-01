package dan

import (
	"log"
	"time"
)

func Logger() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			start := time.Now()

			err := next(c)

			log.Printf("[LOG] %s %s | %v | Error: %v", c.R.Method, c.R.URL.Path, time.Since(start), err)

			return err
		}
	}
}

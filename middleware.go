package dan

import (
	"log"
	"net/http"
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

func Recovery() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					log.Printf("[PANIC] %s %s -> %v", c.R.Method, c.R.URL.Path, recovered)
					err = c.Error(http.StatusInternalServerError, "Internal Server Error")
				}
			}()

			return next(c)
		}
	}
}

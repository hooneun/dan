package dan

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

func Logger() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			start := time.Now()

			err := next(c)

			log.Printf("[LOG] %s %s | %v | Error: %s", c.R.Method, sanitizeLogValue(c.R.URL.Path), time.Since(start), sanitizeLogValue(err))

			return err
		}
	}
}

func Recovery() MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					log.Printf("[PANIC] %s %s -> %s", c.R.Method, sanitizeLogValue(c.R.URL.Path), sanitizeLogValue(recovered))
					if !c.Written() {
						err = c.Error(http.StatusInternalServerError, "Internal Server Error")
					}
				}
			}()

			return next(c)
		}
	}
}

func BodyLimit(maxBytes int64) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Context) error {
			c.R.Body = http.MaxBytesReader(c.W, c.R.Body, maxBytes)
			return next(c)
		}
	}
}

func sanitizeLogValue(value any) string {
	return strings.NewReplacer(
		"\n", "\\n",
		"\r", "\\r",
		"\t", "\\t",
	).Replace(fmt.Sprint(value))
}

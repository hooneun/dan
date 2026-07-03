package dan

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context) error
type MiddlewareFunc func(HandlerFunc) HandlerFunc

type RouterGroup struct {
	prefix      string
	middlewares []MiddlewareFunc
	engine      *Engine
}

func (g *RouterGroup) Group(prefix string, middlewares ...MiddlewareFunc) *RouterGroup {
	combined := make([]MiddlewareFunc, len(g.middlewares)+len(middlewares))
	copy(combined, g.middlewares)
	copy(combined[len(g.middlewares):], middlewares)

	return &RouterGroup{
		prefix:      g.prefix + prefix,
		middlewares: combined,
		engine:      g.engine,
	}
}

func (g *RouterGroup) Use(mw ...MiddlewareFunc) {
	g.middlewares = append(g.middlewares, mw...)
}

func (g *RouterGroup) Handle(method, path string, handler HandlerFunc) {
	fullPath := g.prefix + path
	pattern := method + " " + fullPath

	finalHandler := handler
	for i := len(g.middlewares) - 1; i >= 0; i-- {
		finalHandler = g.middlewares[i](finalHandler)
	}

	g.engine.handlers[pattern] = finalHandler
}

func (g *RouterGroup) GET(path string, h HandlerFunc)     { g.Handle(http.MethodGet, path, h) }
func (g *RouterGroup) POST(path string, h HandlerFunc)    { g.Handle(http.MethodPost, path, h) }
func (g *RouterGroup) PUT(path string, h HandlerFunc)     { g.Handle(http.MethodPut, path, h) }
func (g *RouterGroup) PATCH(path string, h HandlerFunc)   { g.Handle(http.MethodPatch, path, h) }
func (g *RouterGroup) DELETE(path string, h HandlerFunc)  { g.Handle(http.MethodDelete, path, h) }
func (g *RouterGroup) OPTIONS(path string, h HandlerFunc) { g.Handle(http.MethodOptions, path, h) }

type Engine struct {
	*RouterGroup
	handlers map[string]HandlerFunc
}

func NewEngine() *Engine {
	engine := &Engine{handlers: make(map[string]HandlerFunc)}

	engine.RouterGroup = &RouterGroup{
		prefix:      "",
		middlewares: []MiddlewareFunc{},
		engine:      engine,
	}

	engine.middlewares = append(engine.middlewares, Logger())

	return engine
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	pattern := r.Method + " " + r.URL.Path
	if handler, ok := e.handlers[pattern]; ok {
		e.handle(w, r, handler)
		return
	}

	if handler, ok := e.match(r); ok {
		e.handle(w, r, handler)
		return
	}

	http.NotFound(w, r)
}

func (e *Engine) handle(w http.ResponseWriter, r *http.Request, handler HandlerFunc) {
	ctx := &Context{W: w, R: r}
	if err := handler(ctx); err != nil {
		log.Printf("[ERROR] %s %s -> %v", r.Method, r.URL.Path, err)
		ctx.Error(http.StatusInternalServerError, "Internal Server Error")
	}
}

func (e *Engine) match(r *http.Request) (HandlerFunc, bool) {
	for pattern, handler := range e.handlers {
		method, routePath, ok := splitPattern(pattern)
		if !ok || method != r.Method {
			continue
		}

		if matchPath(routePath, r.URL.Path, r) {
			return handler, true
		}
	}

	return nil, false
}

func splitPattern(pattern string) (string, string, bool) {
	parts := strings.SplitN(pattern, " ", 2)
	if len(parts) != 2 {
		return "", "", false
	}

	return parts[0], parts[1], true
}

func matchPath(routePath, requestPath string, r *http.Request) bool {
	routeParts := splitPath(routePath)
	requestParts := splitPath(requestPath)

	if len(routeParts) != len(requestParts) {
		return false
	}

	for i := range routeParts {
		routePart := routeParts[i]
		requestPart := requestParts[i]

		if strings.HasPrefix(routePart, ":") {
			key := strings.TrimPrefix(routePart, ":")
			if key == "" {
				return false
			}

			r.SetPathValue(key, requestPart)
			continue
		}

		if routePart != requestPart {
			return false
		}
	}

	return true
}

func splitPath(path string) []string {
	path = strings.Trim(path, "/")
	if path == "" {
		return []string{}
	}

	return strings.Split(path, "/")
}

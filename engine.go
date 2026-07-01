package dan

import (
	"log"
	"net/http"
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

	g.engine.mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ctx := &Context{W: w, R: r}
		if err := finalHandler(ctx); err != nil {
			log.Printf("[ERROR] %s %s -> %v", r.Method, r.URL.Path, err)
			ctx.Error(http.StatusInternalServerError, "Internal Server Error")
		}
	})
}

func (g *RouterGroup) GET(path string, h HandlerFunc)  { g.Handle(http.MethodGet, path, h) }
func (g *RouterGroup) POST(path string, h HandlerFunc) { g.Handle(http.MethodPost, path, h) }

type Engine struct {
	*RouterGroup
	mux *http.ServeMux
}

func NewEngine() *Engine {
	mux := http.NewServeMux()
	engine := &Engine{mux: mux}

	engine.RouterGroup = &RouterGroup{
		prefix:      "",
		middlewares: []MiddlewareFunc{},
		engine:      engine,
	}

	engine.middlewares = append(engine.middlewares, Logger())

	return engine
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.mux.ServeHTTP(w, r)
}

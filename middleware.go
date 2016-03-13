package middleware

type Handler interface {
	Call(env map[string]interface{})
}

type HandlerFunc func(env map[string]interface{})

func (h HandlerFunc) Call(env map[string]interface{}) {
	h(env)
}

// A handler that responds to a arbitrary request
// Call should write state changes to the map
type MiddlewareHandler interface {
	Call(env map[string]interface{}, next HandlerFunc)
}

// The MiddlewareHandlerFunc type is an adapter to allow the use of ordinary
// functions as MiddlewareHandler handlers. If f is a function with the
// appropriate signature, MiddlewareHandlerFunc(f) is a MiddlewareHandler that calls f.
type MiddlewareHandlerFunc func(env map[string]interface{}, next HandlerFunc)

func (f MiddlewareHandlerFunc) Call(env map[string]interface{}, next HandlerFunc) {
	f(env, next)
}

type link struct {
	handler MiddlewareHandler
	next    *link
}

func (l link) Call(env map[string]interface{}) {
	l.handler.Call(env, l.next.Call)
}

// Chain acts as a list of MiddlewareHandler constructors.
// Chain is effectively immutable:
// once created, it will always hold
// the same set of constructors in the same order.
type Chain struct {
	handlers []MiddlewareHandler
	links    link
}

// New creates a new chain,
// memorizing the given list of middleware constructors.
// New serves no other function,
// constructors are only called upon a call to Then().
func New(handlers ...MiddlewareHandler) *Chain {
	return &Chain{
		handlers: handlers,
		links:    build(handlers),
	}
}

func (c *Chain) Call(env map[string]interface{}) {
	c.links.Call(env)
}

func (c *Chain) Use(handler MiddlewareHandler) {
	c.handlers = append(c.handlers, handler)
	c.links = build(c.handlers)
}

func (c *Chain) UseFunc(h func(env map[string]interface{}, next HandlerFunc)) {
	c.Use(MiddlewareHandlerFunc(h))
}

func build(handlers []MiddlewareHandler) link {
	var next link

	if len(handlers) == 0 {
		return defaultMiddlewareLink
	} else if len(handlers) > 1 {
		next = build(handlers[1:])
	} else {
		next = defaultMiddlewareLink
	}

	return link{handlers[0], &next}
}

type NoopMiddlewareHandler struct {
}

func (h NoopMiddlewareHandler) Call(env map[string]interface{}, next HandlerFunc) {
	// noop
}

var defaultMiddlewareLink = link{NoopMiddlewareHandler{}, &link{}}

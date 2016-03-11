package middleware

// A handler that responds to a arbitrary request
// Call should write state changes to the map
type MiddlewareHandler interface {
	Call(map[string]interface{})
}

type NoopMiddlewareHandler struct {
}

func (h NoopMiddlewareHandler) Call(env map[string]interface{}) {
	// noop
}

var DefaultMiddlewareHandler = NoopMiddlewareHandler{}

// The MiddlewareHandlerFunc type is an adapter to allow the use of ordinary
// functions as MiddlewareHandler handlers. If f is a function with the
// appropriate signature, MiddlewareHandlerFunc(f) is a MiddlewareHandler that calls f.
type MiddlewareHandlerFunc func(map[string]interface{})

func (f MiddlewareHandlerFunc) Call(env map[string]interface{}) {
	f(env)
}

// A constructor for a piece of middleware.
// Some middleware use this constructor out of the box,
// so in most cases you can just pass somepackage.New
type Constructor func(MiddlewareHandler) MiddlewareHandler

// Chain acts as a list of MiddlewareHandler constructors.
// Chain is effectively immutable:
// once created, it will always hold
// the same set of constructors in the same order.
type Chain struct {
	constructors []Constructor
}

// New creates a new chain,
// memorizing the given list of middleware constructors.
// New serves no other function,
// constructors are only called upon a call to Then().
func New(constructors ...Constructor) Chain {
	return Chain{append(([]Constructor)(nil), constructors...)}
}

// Then chains the middleware and returns the final MiddlewareHandler.
//     New(m1, m2, m3).Then(h)
// is equivalent to:
//     m1(m2(m3(h)))
// When the request comes in, it will be passed to m1, then m2, then m3
// and finally, the given handler
// (assuming every middleware calls the following one).
func (c Chain) Then(f MiddlewareHandler) MiddlewareHandler {
	if f == nil {
		f = DefaultMiddlewareHandler
	}

	for i := len(c.constructors) - 1; i >= 0; i-- {
		f = c.constructors[i](f)
	}

	return f
}

// ThenFunc works identically to Then, but takes
// a MiddlewareHandlerFunc instead of a MiddlewareHandler.
//
// The following two statements are equivalent:
//     c.Then(MiddlewareHandlerFunc(fn))
//     c.ThenFunc(fn)
//
// ThenFunc provides all the guarantees of Then.
func (c Chain) ThenFunc(fn MiddlewareHandlerFunc) MiddlewareHandler {
	if fn == nil {
		return c.Then(nil)
	}
	return c.Then(MiddlewareHandlerFunc(fn))
}

// Append extends a chain, adding the specified constructors
// as the last ones in the request flow.
//
// Append returns a new chain, leaving the original one untouched.
//
//     stdChain := middleware.New(m1, m2)
//     extChain := stdChain.Append(m3, m4)
//     // requests in stdChain go m1 -> m2
//     // requests in extChain go m1 -> m2 -> m3 -> m4
func (c Chain) Append(constructors ...Constructor) Chain {
	newCons := make([]Constructor, len(c.constructors)+len(constructors))
	copy(newCons, c.constructors)
	copy(newCons[len(c.constructors):], constructors)

	return New(newCons...)
}

// Extend extends a chain by adding the specified chain
// as the last one in the request flow.
func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.constructors...)
}

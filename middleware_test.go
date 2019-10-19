package middleware

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Not recommended (https://golang.org/pkg/reflect/#Value.Pointer),
// but the best we can do.
func funcsEqual(f1, f2 interface{}) bool {
	val1 := reflect.ValueOf(f1)
	val2 := reflect.ValueOf(f2)
	return val1.Pointer() == val2.Pointer()
}

var testApp = MiddlewareHandlerFunc(func(env map[string]interface{}, next HandlerFunc) {
	fmt.Println("")
})

func TestNewTreatsNilAsEmpty(t *testing.T) {
	chain := New()
	assert.True(t, chain.links == defaultMiddlewareLink)
}

func TestNewTakesFuncs(t *testing.T) {
	c1 := func(env map[string]interface{}, h HandlerFunc) {}
	c2 := func(env map[string]interface{}, h HandlerFunc) {}

	chain := New(c1, c2)

	assert.True(t, funcsEqual(chain.handlers[0], c1))
	assert.True(t, funcsEqual(chain.handlers[1], c2))
}

func TestCallingChainAfterUsingNew(t *testing.T) {
	c1 := func(env map[string]interface{}, h HandlerFunc) {
		env["result"] = "hello"
		h(env)
	}
	c2 := func(env map[string]interface{}, h HandlerFunc) {
		env["fileSize"] = 1000
		h(env)
	}

	chain := New(c1, c2)

	env := make(map[string]interface{})
	chain.Call(env)

	assert.True(t, env["result"] == "hello")
	assert.True(t, env["fileSize"] == 1000)
}

func TestUse(t *testing.T) {
	c1 := MiddlewareHandlerFunc(func(env map[string]interface{}, h HandlerFunc) {})
	c2 := MiddlewareHandlerFunc(func(env map[string]interface{}, h HandlerFunc) {})

	slice := []MiddlewareHandler{c1, c2}

	chain := New()
	chain.Use(slice...)
	assert.True(t, funcsEqual(chain.handlers[0], slice[0]))
	assert.True(t, funcsEqual(chain.handlers[1], slice[1]))
}

func TestUseFunc(t *testing.T) {
	c1 := func(env map[string]interface{}, h HandlerFunc) {}
	c2 := func(env map[string]interface{}, h HandlerFunc) {}

	chain := New()
	chain.UseFunc(c1, c2)
	assert.True(t, funcsEqual(chain.handlers[0], MiddlewareHandlerFunc(c1)))
	assert.True(t, funcsEqual(chain.handlers[1], MiddlewareHandlerFunc(c2)))
}

func TestCallingChain(t *testing.T) {
	c1 := func(env map[string]interface{}, h HandlerFunc) {
		env["result"] = "hello"
		h(env)
	}
	c2 := func(env map[string]interface{}, h HandlerFunc) {
		env["fileSize"] = 1000
		h(env)
	}

	chain := New()
	chain.UseFunc(c1, c2)

	env := make(map[string]interface{})
	chain.Call(env)

	assert.True(t, env["result"] == "hello")
	assert.True(t, env["fileSize"] == 1000)
}

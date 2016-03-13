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

var testApp = MiddlewareHandlerFunc(func(env map[string]interface{}, next MiddlewareHandlerFunc) {
	fmt.Println("")
})

func TestNewTreatsNilAsEmpty(t *testing.T) {
	chain := New()
	assert.True(t, chain.links == defaultMiddlewareLink)
}

func TestUse(t *testing.T) {
	c1 := func(env map[string]interface{}, h MiddlewareHandlerFunc) {}
	c2 := func(env map[string]interface{}, h MiddlewareHandlerFunc) {}

	slice := []MiddlewareHandlerFunc{c1, c2}

	chain := New()
	chain.Use(slice[0])
	chain.Use(slice[1])
	assert.True(t, funcsEqual(chain.handlers[0], slice[0]))
	assert.True(t, funcsEqual(chain.handlers[1], slice[1]))
}

func TestUseFunc(t *testing.T) {
	c1 := func(env map[string]interface{}, h MiddlewareHandlerFunc) {}
	c2 := func(env map[string]interface{}, h MiddlewareHandlerFunc) {}

	chain := New()
	chain.UseFunc(c1)
	chain.UseFunc(c2)
	assert.True(t, funcsEqual(chain.handlers[0], MiddlewareHandlerFunc(c1)))
	assert.True(t, funcsEqual(chain.handlers[1], MiddlewareHandlerFunc(c2)))
}
